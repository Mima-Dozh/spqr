package relay

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgproto3"
	"github.com/pg-sharding/lyx/lyx"
	"github.com/pg-sharding/spqr/pkg/config"
	"github.com/pg-sharding/spqr/pkg/models/kr"
	"github.com/pg-sharding/spqr/pkg/models/spqrerror"
	"github.com/pg-sharding/spqr/pkg/prepstatement"
	"github.com/pg-sharding/spqr/pkg/session"
	"github.com/pg-sharding/spqr/pkg/spqrlog"
	"github.com/pg-sharding/spqr/pkg/txstatus"
	"github.com/pg-sharding/spqr/router/parser"
	"github.com/pg-sharding/spqr/router/routehint"
	"github.com/pg-sharding/spqr/router/routingstate"
	"github.com/pg-sharding/spqr/router/statistics"
)

func AdvancedPoolModeNeeded(rst RelayStateMgr) bool {
	return rst.Client().Rule().PoolMode == config.PoolModeTransaction && rst.Client().Rule().PoolPreparedStatement || rst.RouterMode() == config.ProxyMode
}

func deparseRouteHint(rst RelayStateMgr, params map[string]string) (routehint.RouteHint, error) {
	if _, ok := params[session.SPQR_SCATTER_QUERY]; ok {
		return &routehint.ScatterRouteHint{}, nil
	}
	if val, ok := params[session.SPQR_SHARDING_KEY]; ok {
		spqrlog.Zero.Debug().Str("sharding key", val).Msg("checking hint key")

		dsId := ""
		if dsId, ok = params[session.SPQR_DISTRIBUTION]; !ok {
			return nil, spqrerror.New(spqrerror.SPQR_NO_DISTRIBUTION, "sharding key in comment without distribution")
		}

		ctx := context.TODO()
		krs, err := rst.QueryRouter().Mgr().ListKeyRanges(ctx, dsId)
		if err != nil {
			return nil, err
		}

		distrib, err := rst.QueryRouter().Mgr().GetDistribution(ctx, dsId)
		if err != nil {
			return nil, err
		}

		// TODO: fix this
		compositeKey, err := kr.KeyRangeBoundFromStrings(distrib.ColTypes, []string{val})

		if err != nil {
			return nil, err
		}

		ds, err := rst.QueryRouter().DeparseKeyWithRangesInternal(ctx, compositeKey, krs)
		if err != nil {
			return nil, err
		}
		return &routehint.TargetRouteHint{
			State: routingstate.ShardMatchState{
				Route: ds,
			},
		}, nil
	}

	return &routehint.EmptyRouteHint{}, nil
}

// ProcQueryAdvanced processes query, with router relay state
// There are several types of query that we want to process in non-passthrough way.
// For example, after BEGIN we wait until first client query witch can be router to some shard.
// So, we need to proccess SETs, BEGINs, ROLLBACKs etc ourselves.
// ProtoStateHandler provides set of function for either simple of extended protoc interactions
// query param is either plain query from simple proto or bind query from x proto
func ProcQueryAdvanced(rst RelayStateMgr, query string, ph ProtoStateHandler, binderQ func() error, doCaching bool) error {
	statistics.RecordStartTime(statistics.Router, time.Now(), rst.Client().ID())

	spqrlog.Zero.Debug().Str("query", query).Uint("client", spqrlog.GetPointer(rst.Client())).Msgf("process relay state advanced")
	state, comment, err := rst.Parse(query, doCaching)
	if err != nil {
		return fmt.Errorf("error processing query '%v': %w", query, err)
	}

	queryProc := func() error {
		mp, err := parser.ParseComment(comment)

		if err == nil {
			if val, ok := mp["target-session-attrs"]; ok {
				// TBD: validate
				spqrlog.Zero.Debug().Str("tsa", val).Msg("parse tsa from comment")
				rst.Client().SetTsa(val)
			}
			recheckRouteHint := false

			if val, ok := mp[session.SPQR_DEFAULT_ROUTE_BEHAVIOUR]; ok {
				recheckRouteHint = true
				spqrlog.Zero.Debug().Str("default route", val).Msg("parse default route behaviour from comment")
				rst.Client().SetDefaultRouteBehaviour(val)
			}

			val, ok := mp[session.SPQR_SHARDING_KEY]
			recheckRouteHint = recheckRouteHint || ok || rst.Client().ShardingKey() != ""
			if ok {
				spqrlog.Zero.Debug().Str("sharding key", val).Msg("parse sharding key from comment")
				rst.Client().SetShardingKey(val)
			}

			if recheckRouteHint {
				/* if not enforsed by current query, use param, if set */
				if rst.Client().ShardingKey() != "" {
					mp[session.SPQR_SHARDING_KEY] = rst.Client().ShardingKey()
				}
				routeHint, err := deparseRouteHint(rst, mp)
				if err == nil {
					spqrlog.Zero.Debug().Interface("hint", routeHint).Msg("setting routing hint")
					rst.Client().SetRouteHint(routeHint)
				} else {
					spqrlog.Zero.Debug().Err(err).Msg("failed to deparse routing hint")
				}
			}

			if val, ok := mp[session.SPQR_AUTO_DISTRIBUTION]; ok {
				if valDistrib, ok := mp[session.SPQR_DISTRIBUTION_KEY]; ok {
					_, err = rst.QueryRouter().Mgr().GetDistribution(context.TODO(), val)
					if err != nil {
						return err
					}

					/* This is an ddl query, which creates relation along with attaching to dsitribution */
					rst.Client().SetAutoDistribution(val)
					rst.Client().SetDistributionKey(valDistrib)

					/* this is too early to do anything with distribution hint, as we do not yet parsed
					* DDL of about-to-be-created relation
					 */
				} else {
					return fmt.Errorf("spqr distribution specified, but distribution key omitted.")
				}
			}
			if val, ok := mp[session.SPQR_ALLOW_MULTISHARD]; ok && val == "true" {
				rst.Client().SetAllowMultishard(true)
			} else {
				rst.Client().SetAllowMultishard(false)
			}

			if val, ok := mp[session.SPQR_EXECUTE_ON]; ok {
				if _, ok := config.RouterConfig().ShardMapping[val]; !ok {
					return fmt.Errorf("no such shard: %v", val)
				}
				rst.Client().SetExecuteOn(val)
			} else {
				/* unset or reset if any */
				rst.Client().SetExecuteOn("")
			}
		}

		return binderQ()
	}

	switch st := state.(type) {
	case parser.ParseStateTXBegin:
		if rst.TxStatus() != txstatus.TXIDLE {
			// ignore this
			_ = rst.Client().ReplyWarningf("there is already transaction in progress")
			return rst.Client().ReplyCommandComplete("BEGIN")
		}
		// explicitly set silent query message, as it can differ from query begin in xporot
		rst.AddSilentQuery(&pgproto3.Query{
			String: query,
		})

		rst.SetTxStatus(txstatus.TXACT)
		rst.Client().StartTx()

		spqrlog.Zero.Debug().Msg("start new transaction")

		for _, opt := range st.Options {
			switch opt {
			case lyx.TransactionReadOnly:
				rst.Client().SetTsa(config.TargetSessionAttrsPS)
			case lyx.TransactionReadWrite:
				rst.Client().SetTsa(config.TargetSessionAttrsRW)
			}
		}
		return rst.Client().ReplyCommandComplete("BEGIN")
	case parser.ParseStateTXCommit:
		if rst.TxStatus() != txstatus.TXACT && rst.TxStatus() != txstatus.TXERR {
			_ = rst.Client().ReplyWarningf("there is no transaction in progress")
			return rst.Client().ReplyCommandComplete("COMMIT")
		}
		return ph.ExecCommit(rst, query)
	case parser.ParseStateTXRollback:
		if rst.TxStatus() != txstatus.TXACT && rst.TxStatus() != txstatus.TXERR {
			_ = rst.Client().ReplyWarningf("there is no transaction in progress")
			return rst.Client().ReplyCommandComplete("ROLLBACK")
		}
		return ph.ExecRollback(rst, query)
	case parser.ParseStateEmptyQuery:
		if err := rst.Client().Send(&pgproto3.EmptyQueryResponse{}); err != nil {
			return err
		}
		// do not complete relay  here
		return nil
	// with tx pooling we might have no active connection while processing set x to y
	case parser.ParseStateSetStmt:
		spqrlog.Zero.Debug().
			Str("name", st.Name).
			Str("value", st.Value).
			Msg("applying parsed set stmt")

		if strings.HasPrefix(st.Name, "__spqr__") {
			switch st.Name {
			case session.SPQR_DISTRIBUTION:
				rst.Client().SetDistribution(st.Value)
			case session.SPQR_DEFAULT_ROUTE_BEHAVIOUR:
				rst.Client().SetDefaultRouteBehaviour(st.Value)
			case session.SPQR_SHARDING_KEY:
				rst.Client().SetShardingKey(st.Value)
			case session.SPQR_REPLY_NOTICE:
				if st.Value == "on" || st.Value == "true" {
					rst.Client().SetShowNoticeMsg(true)
				} else {
					rst.Client().SetShowNoticeMsg(false)
				}
			case session.SPQR_MAINTAIN_PARAMS:
				if st.Value == "on" || st.Value == "true" {
					rst.Client().SetMaintainParams(true)
				} else {
					rst.Client().SetMaintainParams(false)
				}
			default:
				rst.Client().SetParam(st.Name, st.Value)
			}

			routeHint, err := deparseRouteHint(rst, map[string]string{
				session.SPQR_DISTRIBUTION: rst.Client().Distribution(),
				session.SPQR_SHARDING_KEY: rst.Client().ShardingKey(),
			})

			if err == nil {
				spqrlog.Zero.Debug().Interface("hint", routeHint).Msg("setting routing hint")
				rst.Client().SetRouteHint(routeHint)
			} else {
				spqrlog.Zero.Debug().Err(err).Msg("failed to deparse routing hint")
			}

			return rst.Client().ReplyCommandComplete("SET")
		}

		return ph.ExecSet(rst, query, st.Name, st.Value)
	case parser.ParseStateShowStmt:
		param := st.Name
		// manually create router responce
		// here we just reply single row with single column value

		switch param {
		case session.SPQR_DISTRIBUTION:
			_ = rst.Client().Send(&pgproto3.ErrorResponse{
				Message: fmt.Sprintf("parameter \"%s\" isn't user accessible", session.SPQR_DISTRIBUTION),
			})
		case session.SPQR_DEFAULT_ROUTE_BEHAVIOUR:

			_ = rst.Client().Send(
				&pgproto3.RowDescription{
					Fields: []pgproto3.FieldDescription{
						{
							Name:         []byte("default route behaviour"),
							DataTypeOID:  25,
							DataTypeSize: -1,
							TypeModifier: -1,
						},
					},
				},
			)

			_ = rst.Client().Send(
				&pgproto3.DataRow{
					Values: [][]byte{
						[]byte(rst.Client().DefaultRouteBehaviour()),
					},
				},
			)

		case session.SPQR_REPLY_NOTICE:

			_ = rst.Client().Send(
				&pgproto3.RowDescription{
					Fields: []pgproto3.FieldDescription{
						{
							Name:         []byte("show notice messages"),
							DataTypeOID:  25,
							DataTypeSize: -1,
							TypeModifier: -1,
						},
					},
				},
			)

			if rst.Client().ShowNoticeMsg() {
				_ = rst.Client().Send(
					&pgproto3.DataRow{
						Values: [][]byte{
							[]byte("true"),
						},
					},
				)
			} else {
				_ = rst.Client().Send(
					&pgproto3.DataRow{
						Values: [][]byte{
							[]byte("false"),
						},
					},
				)
			}

		case session.SPQR_MAINTAIN_PARAMS:

			_ = rst.Client().Send(
				&pgproto3.RowDescription{
					Fields: []pgproto3.FieldDescription{
						{
							Name:         []byte("maintain params"),
							DataTypeOID:  25,
							DataTypeSize: -1,
							TypeModifier: -1,
						},
					},
				},
			)

			if rst.Client().MaintainParams() {
				_ = rst.Client().Send(
					&pgproto3.DataRow{
						Values: [][]byte{
							[]byte("true"),
						},
					},
				)
			} else {
				_ = rst.Client().Send(
					&pgproto3.DataRow{
						Values: [][]byte{
							[]byte("false"),
						},
					},
				)
			}

		case session.SPQR_SHARDING_KEY:

			_ = rst.Client().Send(
				&pgproto3.RowDescription{
					Fields: []pgproto3.FieldDescription{
						{
							Name:         []byte("sharding key"),
							DataTypeOID:  25,
							DataTypeSize: -1,
							TypeModifier: -1,
						},
					},
				},
			)

			_ = rst.Client().Send(
				&pgproto3.DataRow{
					Values: [][]byte{
						[]byte("no val"),
					},
				},
			)
		case session.SPQR_SCATTER_QUERY:

			_ = rst.Client().Send(
				&pgproto3.RowDescription{
					Fields: []pgproto3.FieldDescription{
						{
							Name:         []byte("scatter query"),
							DataTypeOID:  25,
							DataTypeSize: -1,
							TypeModifier: -1,
						},
					},
				},
			)

			_ = rst.Client().Send(
				&pgproto3.DataRow{
					Values: [][]byte{
						[]byte("no val"),
					},
				},
			)
		default:

			/* If router does dot have any info about param, fire query to random shard. */
			if _, ok := rst.Client().Params()[param]; !ok {
				return queryProc()
			}

			_ = rst.Client().Send(
				&pgproto3.RowDescription{
					Fields: []pgproto3.FieldDescription{
						{
							Name:         []byte(param),
							DataTypeOID:  25,
							DataTypeSize: -1,
							TypeModifier: -1,
						},
					},
				},
			)
			_ = rst.Client().Send(
				&pgproto3.DataRow{
					Values: [][]byte{
						[]byte(rst.Client().Params()[param]),
					},
				},
			)
		}
		return rst.Client().ReplyCommandComplete("SHOW")
	case parser.ParseStateResetStmt:
		rst.Client().ResetParam(st.Name)

		if err := ph.ExecReset(rst, query, st.Name); err != nil {
			return err
		}

		return rst.Client().ReplyCommandComplete("RESET")
	case parser.ParseStateResetMetadataStmt:
		if err := ph.ExecResetMetadata(rst, query, st.Setting); err != nil {
			return err
		}

		rst.Client().ResetParam(st.Setting)
		if st.Setting == "session_authorization" {
			rst.Client().ResetParam("role")
		}

		return rst.Client().ReplyCommandComplete("RESET")
	case parser.ParseStateResetAllStmt:
		rst.Client().ResetAll()

		return rst.Client().ReplyCommandComplete("RESET")
	case parser.ParseStateSetLocalStmt:
		if err := ph.ExecSetLocal(rst, query, st.Name, st.Value); err != nil {
			return err
		}

		return rst.Client().ReplyCommandComplete("SET")
	case parser.ParseStatePrepareStmt:
		// sql level prepares stmt pooling
		if AdvancedPoolModeNeeded(rst) {
			spqrlog.Zero.Debug().Msg("sql level prep statement pooling support is on")

			/* no OIDS for SQL level prep stmt */
			rst.Client().StorePreparedStatement(&prepstatement.PreparedStatementDefinition{
				Name:  st.Name,
				Query: st.Query,
			})
			return nil
		} else {
			// process like regular query
			return queryProc()
		}
	case parser.ParseStateExecute:
		if AdvancedPoolModeNeeded(rst) {
			// do nothing
			// wtf? TODO: test and fix
			rst.Client().PreparedStatementQueryByName(st.Name)
			return nil
		} else {
			// process like regular query
			return queryProc()
		}
	case parser.ParseStateExplain:
		_ = rst.Client().ReplyErrMsgByCode(spqrerror.SPQR_UNEXPECTED)
		return nil
	default:
		return queryProc()
	}
}
