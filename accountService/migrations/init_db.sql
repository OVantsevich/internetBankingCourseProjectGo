create table users
(
    id                serial
        constraint users_pk
            primary key,
    user_login        varchar(50)                               not null,
    user_name         varchar(50)                               not null,
    surname           varchar(50)                               not null,
    is_deleted        boolean      default false                not null
);

alter table users
    owner to postgres;

create unique index users_id_uindex
    on users (id);

create unique index users_user_login_uindex
    on users (user_login);

create table accounts
(
    id                serial
        constraint account_pk
            primary key,
    user_id           integer                                           not null
        constraint accounts_user_id_fk
            references users,
    account_name      varchar(40)  default 'account'::character varying not null,
    amount            integer      default 0                            not null,
    is_deleted        boolean      default false                        not null,
    creation_date     timestamp(6) default CURRENT_TIMESTAMP(6)         not null,
    modification_date timestamp(6) default CURRENT_TIMESTAMP(6)         not null
);

alter table accounts
    owner to postgres;

create unique index account_id_uindex
    on accounts (id);

create index accounts_is_deleted_index
    on accounts (is_deleted);

create unique index account_user_id_account_name_uindex
    on accounts (user_id, account_name);

create table account_types
(
    id                serial
        constraint account_types_pk
            primary key,
    account_type_name varchar(40)  default 'type'::character varying not null
        constraint account_types_account_type_name_unq
            unique,
    is_deleted        boolean      default false                     not null,
    creation_date     timestamp(6) default CURRENT_TIMESTAMP(6)      not null,
    modification_date timestamp(6) default CURRENT_TIMESTAMP(6)      not null
);

alter table account_types
    owner to postgres;

create unique index account_types_uindex
    on account_types (id);

create unique index account_types_account_type_name_uindex
    on account_types (account_type_name);

create index account_types_is_deleted_index
    on account_types (is_deleted);

insert into account_types (id, account_type_name)
values (1, 'dollar_account');

insert into account_types (id, account_type_name)
values (2, 'euro_account');

insert into account_types (id, account_type_name)
values (3, 'ruble_account');

create table l_account_types_accounts
(
    id                serial
        constraint l_account_types_accounts_pk
            primary key,
    account_type_id   integer                                   not null
        constraint l_account_types_accounts_account_type_id_fk
            references account_types,
    account_id        integer                                   not null
        constraint l_account_types_accounts_account_id_fk
            references accounts,
    is_deleted        boolean      default false                not null,
    creation_date     timestamp(6) default CURRENT_TIMESTAMP(6) not null,
    modification_date timestamp(6) default CURRENT_TIMESTAMP(6) not null
);

alter table l_account_types_accounts
    owner to postgres;

create unique index l_account_types_accounts_account_type_id_account_id_uindex
    on l_account_types_accounts (account_type_id, account_id);

create index l_account_types_accounts_is_deleted_index
    on account_types (is_deleted);

create table transactions
(
    id                serial
        constraint transaction_pk
            primary key,
    account_id        integer                                   not null
        constraint transactions_account_id_fk
            references accounts,
    amount            integer      default 0                    not null,
    is_deleted        boolean      default false                not null,
    creation_date     timestamp(6) default CURRENT_TIMESTAMP(6) not null,
    modification_date timestamp(6) default CURRENT_TIMESTAMP(6) not null
);

alter table transactions
    owner to postgres;

create unique index transactions_id_uindex
    on transactions (id);

create index transactions_account_id_index
    on transactions (account_id);

create index transactions_is_deleted_index
    on transactions (is_deleted);

create table transaction_types
(
    id                    serial
        constraint transaction_types_pk
            primary key,
    transaction_type_name varchar(40)  default 'type'::character varying not null
        constraint transaction_types_transaction_type_name_unq
            unique,
    is_deleted            boolean      default false                     not null,
    creation_date         timestamp(6) default CURRENT_TIMESTAMP(6)      not null,
    modification_date     timestamp(6) default CURRENT_TIMESTAMP(6)      not null
);

alter table transaction_types
    owner to postgres;

create unique index transaction_types_uindex
    on transaction_types (id);

create unique index transaction_types_transaction_type_name_uindex
    on transaction_types (transaction_type_name);

create index transaction_types_is_deleted_index
    on transaction_types (is_deleted);

insert into transaction_types (id, transaction_type_name)
values (1, 'internal_transaction');

insert into transaction_types (id, transaction_type_name)
values (2, 'external_transaction');

create table l_transaction_types_transactions
(
    id                  serial
        constraint l_transaction_types_transactions_pk
            primary key,
    transaction_type_id integer                                   not null
        constraint l_transaction_types_transactions_transaction_type_id_fk
            references transaction_types,
    transaction_id      integer                                   not null
        constraint l_transaction_types_transactions_transaction_id_fk
            references transactions,
    is_deleted          boolean      default false                not null,
    creation_date       timestamp(6) default CURRENT_TIMESTAMP(6) not null,
    modification_date   timestamp(6) default CURRENT_TIMESTAMP(6) not null
);

alter table l_transaction_types_transactions
    owner to postgres;

create unique index l_transaction_types_transactions_type_id_id_uindex
    on l_transaction_types_transactions (transaction_type_id, transaction_id);

create index l_transaction_types_transactions_is_deleted_index
    on l_transaction_types_transactions (is_deleted);

CREATE OR REPLACE FUNCTION fun_updateaccountamount() RETURNS TRIGGER AS
$BODY$
BEGIN
    update accounts set amount = amount - new.amount
    where new.account_id = id;
    RETURN new;
END;
$BODY$
    language plpgsql;

CREATE TRIGGER TRI_TRANSACTIONS
    AFTER INSERT ON TRANSACTIONS
    FOR EACH ROW
EXECUTE PROCEDURE fun_updateaccountamount();