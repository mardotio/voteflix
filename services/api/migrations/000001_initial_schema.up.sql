create extension if not exists "uuid-ossp";

create or replace function gen_random_movie_seed()
returns integer
as
$$
    begin
        return floor(random() * (2147483647::bigint - -2147483648::bigint) - 2147483648::bigint)::integer;
    end;
$$ language 'plpgsql' STRICT;

create table if not exists users (
    id uuid primary key default gen_random_uuid(),
    discord_id varchar(100) not null unique,
    discord_username varchar(100) not null,
    discord_avatar_id varchar(100),
    created_at timestamp default current_timestamp,
    updated_at timestamp
);

create table if not exists lists (
    id uuid primary key default gen_random_uuid(),
    name varchar(100) not null,
    discord_server_id varchar(100) not null unique,
    creator_id uuid not null,
    created_at timestamp default current_timestamp,
    updated_at timestamp,
    discord_avatar_id varchar(100),

    foreign key (creator_id) references users(id)
);

create table if not exists list_users (
    user_id uuid not null,
    list_id uuid not null,
    discord_nickname varchar(32),
    created_at timestamp default current_timestamp,
    updated_at timestamp,

    foreign key (user_id) references users(id),
    foreign key (list_id) references lists(id),
    primary key (user_id, list_id)
);

create table if not exists movies (
    id uuid primary key default gen_random_uuid(),
    list_id uuid not null,
    name varchar(255) not null,
    status varchar(15) not null,
    approve_count integer default 0,
    reject_count integer default 0,
    seed integer default gen_random_movie_seed(),
    creator_id uuid not null,
    created_at timestamp default current_timestamp,
    updated_at timestamp,
    watched_at timestamp,

    foreign key (list_id) references lists(id),
    foreign key (creator_id) references users(id),

    constraint movies_valid_status check (status in ('pending', 'approved', 'watched', 'rejected')),
    constraint movies_valid_approve_count check (approve_count > -1),
    constraint movies_valid_reject_count check (reject_count > -1),
    constraint movies_valid_status_and_watched_timestamp check ((status = 'watched' and watched_at is not null) or (status != 'watched' and watched_at is null))
);

create table if not exists ratings (
    user_id uuid not null,
    movie_id uuid not null,
    rating smallint not null,
    created_at timestamp default current_timestamp,
    updated_at timestamp,

    foreign key (user_id) references users(id),
    foreign key (movie_id) references movies(id),
    primary key (movie_id, user_id)
);

create table if not exists votes (
    user_id uuid not null,
    movie_id uuid not null,
    is_approval boolean not null,
    created_at timestamp default current_timestamp,
    updated_at timestamp,

    foreign key (user_id) references users(id),
    foreign key (movie_id) references movies(id),
    primary key (user_id, movie_id)
);