create table users
(
    id                serial
        constraint users_pk
            primary key,
    user_login        varchar(50)                               not null,
    user_email        varchar(50)                               not null,
    user_password     varchar(200)                              not null,
    user_name         varchar(50)                               not null,
    surname           varchar(50)                               not null,
    is_deleted        boolean      default false                not null,
    creation_date     timestamp(6) default CURRENT_TIMESTAMP(6) not null,
    modification_date timestamp(6) default CURRENT_TIMESTAMP(6) not null
);

alter table users
    owner to postgres;

create unique index users_id_uindex
    on users (id);

create unique index users_user_login_uindex
    on users (user_login);

create unique index users_user_email_uindex
    on users (user_email);
