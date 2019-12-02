create table stats
(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    timestamp DATETIME not null,
    data TEXT not null
);

create index stats_timestamp_index
    on stats (timestamp);

