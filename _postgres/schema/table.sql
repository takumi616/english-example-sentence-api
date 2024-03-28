create table sentence(
    id       serial     PRIMARY KEY,
    body    varchar(500)   not null,
    vocabularies varchar(30)[3]  not null,
    created  timestamp  not null,
    updated  timestamp  not null
);



