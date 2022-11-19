create table users
(
    id                serial
        constraint users_pk
            primary key,
    user_login        varchar(100)                              not null
        constraint constraint_name
            unique,
    user_password     varchar(200)                              not null,
    user_name         varchar(20)                               not null,
    surname           varchar(50)                               not null,
    is_deleted        boolean      default false                not null,
    creation_date     timestamp(6) default CURRENT_TIMESTAMP(6) not null,
    modification_date timestamp(6) default CURRENT_TIMESTAMP(6) not null
);

alter table users
    owner to postgres;

create unique index users_id_uindex
    on users (id);

