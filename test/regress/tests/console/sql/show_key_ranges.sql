CREATE KEY RANGE krid1 FROM 1 ROUTE TO sh1;
CREATE KEY RANGE krid2 FROM 11 ROUTE TO sh1;

SHOW key_ranges;

DROP DATASPACE ALL CASCADE;
DROP SHARDING RULE ALL;
DROP KEY RANGE ALL;