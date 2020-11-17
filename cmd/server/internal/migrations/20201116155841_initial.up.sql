create table buff (
    id serial primary key,
    question text not null,
    answers text[]
);

create table stream (
    id serial primary  key,
    name text not null
);

create table buff_to_stream (
    stream_id int references stream(id),
    buff_id int references buff(id)
);