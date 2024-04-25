create table sentence(
    id       serial     PRIMARY KEY,
    body    varchar(120)   not null,
    vocabularies varchar(15)[3]  not null,
    created  varchar(60)  not null,
    updated  varchar(60)  not null
);



