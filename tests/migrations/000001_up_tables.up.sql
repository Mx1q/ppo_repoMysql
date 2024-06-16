set names utf8;

create table if not exists saladRecipes.user (
        id varchar(36) default (uuid()) primary key,
        name varchar(64) not null,
        email varchar(64) not null check ( email like '%@%.%' ) unique,
        login varchar(64) not null unique,
        password varchar(256) not null,
        role varchar(25) not null default 'user'
    );

create table if not exists saladRecipes.saladType (
        id varchar(36) default (uuid()) primary key,
        name varchar(25) not null,
        description varchar(256) default ''
    );

create table if not exists saladRecipes.salad (
        id varchar(36) default (uuid()) primary key,
        name varchar(32) not null,
        authorId varchar(36),
        description varchar(64) not null default '',
        foreign key (authorId) references saladRecipes.user(id)
    );

create table if not exists saladRecipes.modStatus (
        id serial primary key,
        name varchar(32),
        description varchar(256)
    );

create table if not exists saladRecipes.recipe (
        id varchar(36) default (uuid()) primary key,
        saladId varchar(36) not null,
        status bigint unsigned not null,
        numberOfServings int not null, check ( numberOfServings > 0 ),
        timeToCook int not null, check ( timeToCook > 0 ),
        rating decimal(3, 1) default 0.0, check ( rating >= 0 ), check ( rating <= 5 ),
        foreign key (saladId) references saladRecipes.salad(id) on delete cascade,
        foreign key (status) references saladRecipes.modStatus(id)
    );

create table if not exists saladRecipes.ingredientType (
        id varchar(36) default (uuid()) primary key,
        name varchar(32) not null,
        description varchar(256) default ''
    );

create table if not exists saladRecipes.ingredient (
        id varchar(36) default (uuid()) primary key,
        name varchar(32) not null,
        calories int not null, check ( calories >= 0 ),
        type varchar(36) not null,
        foreign key (type) references saladRecipes.ingredientType(id)
    );

create table if not exists saladRecipes.comment (
        id varchar(36) default (uuid()) primary key,
        author varchar(36) not null,
        salad varchar(36) not null,
        text varchar(64) default '',
        rating int not null, check ( rating >= 1 ), check ( rating <= 5 ),
        unique (author, salad),
        foreign key (author) references saladRecipes.user(id),
        foreign key (salad) references saladRecipes.salad(id) on delete cascade
    );

create table if not exists saladRecipes.measurement (
        id varchar(36) default (uuid()) primary key,
        name varchar(32) not null,
        grams int, check ( grams > 0 )
    );

create table if not exists saladRecipes.recipeStep (
        id varchar(36) default (uuid()) primary key,
        name varchar(32) not null,
        description text not null,
        recipeId varchar(36) not null,
        stepNum int not null, check ( stepNum > 0 ),
        foreign key (recipeId) references saladRecipes.recipe(id) on delete cascade
    );

#   Links between tables
create table if not exists saladRecipes.recipeIngredient (
        id varchar(36) default (uuid()) primary key,
        recipeId varchar(36) not null,
        ingredientId varchar(36) not null,
        measurement varchar(36) not null default '01000000-0000-0000-0000-000000000000',
        amount int not null default 1, check ( amount > 0 ),
        unique (recipeId, ingredientId),
        foreign key (recipeId) references saladRecipes.recipe(id) on delete cascade,
        foreign key (measurement) references saladRecipes.measurement(id),
        foreign key (ingredientId) references saladRecipes.ingredient(id)
    );

create table if not exists saladRecipes.typesOfSalads (
        id varchar(36) default (uuid()) primary key,
        saladId varchar(36) not null,
        typeId varchar(36) not null,
        foreign key (saladId) references saladRecipes.salad(id) on delete cascade,
        foreign key (typeId) references saladRecipes.saladType(id)
    );

create table if not exists saladRecipes.word (
        id varchar(36) default (uuid()) primary key,
        word varchar(32)
    );


insert into saladRecipes.ingredientType(name)
values
    ('фрукт'),
    ('овощ'),
    ('мясо'),
    ('рыба'),
    ('молоко');

insert into saladRecipes.ingredient(id, name, calories, type)
values ('f1fc4bfc-799c-4471-a971-1bb00f7dd30a', 'яблоко', 1, (select id from saladRecipes.ingredientType where name = 'фрукт')),
       ('01000000-0000-0000-0000-000000000000', 'морковь', 2, (select id from saladRecipes.ingredientType where name = 'овощ')),
       ('02000000-0000-0000-0000-000000000000', 'говядина', 3, (select id from saladRecipes.ingredientType where name = 'мясо')),
       ('03000000-0000-0000-0000-000000000000', 'лосось', 4,  (select id from saladRecipes.ingredientType where name = 'рыба')),
       ('04000000-0000-0000-0000-000000000000', 'молоко', 5, (select id from saladRecipes.ingredientType where name = 'молоко'));

insert into saladRecipes.salad(id, name)
values ('fbabc2aa-cd4a-42b0-b68d-d3cf67fba06f', 'цезарь'),
       ('01000000-0000-0000-0000-000000000000', 'овощной'),
       ('02000000-0000-0000-0000-000000000000', 'сезонный'),
       ('03000000-0000-0000-0000-000000000000', 'сельдь под шубой'),
       ('04000000-0000-0000-0000-000000000000', 'греческий');

insert into saladRecipes.saladType(id, name)
values ('7e17866b-2b97-4d2b-b399-42ceeebd5480', 'зима'),
       ('01000000-0000-0000-0000-000000000000', 'лето');

insert into saladRecipes.saladType(name)
values
    ('осень'),
    ('весна'),
    ('мясной');

insert into saladRecipes.typesOfSalads(saladid, typeid)
values
    ((select id from saladRecipes.salad where name = 'цезарь'),
     (select id from saladRecipes.saladType where name = 'зима')),

    ((select id from saladRecipes.salad where name = 'овощной'),
     (select id from saladRecipes.saladType where name = 'лето')),
    ((select id from saladRecipes.salad where name = 'овощной'),
     (select id from saladRecipes.saladType where name = 'зима')),

    ((select id from saladRecipes.salad where name = 'сезонный'),
     (select id from saladRecipes.saladType where name = 'лето')),
    ((select id from saladRecipes.salad where name = 'сезонный'),
     (select id from saladRecipes.saladType where name = 'зима')),
    ((select id from saladRecipes.salad where name = 'сезонный'),
     (select id from saladRecipes.saladType where name = 'весна')),
    ((select id from saladRecipes.salad where name = 'сезонный'),
     (select id from saladRecipes.saladType where name = 'осень')),

    ((select id from saladRecipes.salad where name = 'сельдь под шубой'),
     (select id from saladRecipes.saladType where name = 'зима')),
    ((select id from saladRecipes.salad where name = 'сельдь под шубой'),
     (select id from saladRecipes.saladType where name = 'мясной')),

    ((select id from saladRecipes.salad where name = 'греческий'),
     (select id from saladRecipes.saladType where name = 'зима'));

insert into saladRecipes.modStatus(name)
values
    ('редактирование'),
    ('на модерации'),
    ('отклонено'),
    ('опубликовано'),
    ('снято с публикации');

insert into saladRecipes.recipe(id, saladid, status, numberofservings, timetocook)
values
    ('01000000-0000-0000-0000-000000000000', (select id from saladRecipes.salad where name = 'цезарь'),
     (select id from saladRecipes.modStatus where name = 'опубликовано'),
     1, 1),
    ('02000000-0000-0000-0000-000000000000', (select id from saladRecipes.salad where name = 'овощной'),
     (select id from saladRecipes.modStatus where name = 'опубликовано'),
     2, 2),
    ('03000000-0000-0000-0000-000000000000', (select id from saladRecipes.salad where name = 'сельдь под шубой'),
     (select id from saladRecipes.modStatus where name = 'опубликовано'),
     3, 3),
    ('04000000-0000-0000-0000-000000000000', (select id from saladRecipes.salad where name = 'сезонный'),
     (select id from saladRecipes.modStatus where name = 'опубликовано'),
     4, 4),
    ('05000000-0000-0000-0000-000000000000', (select id from saladRecipes.salad where name = 'греческий'),
     (select id from saladRecipes.modStatus where name = 'опубликовано'),
     5, 5);

insert into saladRecipes.measurement(name, grams)
values
    ('граммов', 1),
    ('чайная ложка', 1),
    ('штук', 1),
    ('килограмм', 1000);

insert into saladRecipes.recipeIngredient(recipeid, ingredientid, measurement, amount)
values
    ((select id from saladRecipes.recipe where numberofservings = 1),
     (select id from saladRecipes.ingredient where name = 'яблоко'),
     (select id from saladRecipes.measurement where name = 'граммов'),
     1),

    ((select id from saladRecipes.recipe where numberofservings = 2),
     (select id from saladRecipes.ingredient where name = 'морковь'),
     (select id from saladRecipes.measurement where name = 'граммов'),
     2),
    ((select id from saladRecipes.recipe where numberofservings = 2),
     (select id from saladRecipes.ingredient where name = 'яблоко'),
     (select id from saladRecipes.measurement where name = 'граммов'),
     3),

    ((select id from saladRecipes.recipe where numberofservings = 3),
     (select id from saladRecipes.ingredient where name = 'говядина'),
     (select id from saladRecipes.measurement where name = 'граммов'),
     4),
    ((select id from saladRecipes.recipe where numberofservings = 3),
     (select id from saladRecipes.ingredient where name = 'яблоко'),
     (select id from saladRecipes.measurement where name = 'граммов'),
     5),

    ((select id from saladRecipes.recipe where numberofservings = 4),
     (select id from saladRecipes.ingredient where name = 'говядина'),
     (select id from saladRecipes.measurement where name = 'граммов'),
     6),
    ((select id from saladRecipes.recipe where numberofservings = 4),
     (select id from saladRecipes.ingredient where name = 'яблоко'),
     (select id from saladRecipes.measurement where name = 'граммов'),
     7),
    ((select id from saladRecipes.recipe where numberofservings = 4),
     (select id from saladRecipes.ingredient where name = 'морковь'),
     (select id from saladRecipes.measurement where name = 'граммов'),
     8),
    ((select id from saladRecipes.recipe where numberofservings = 4),
     (select id from saladRecipes.ingredient where name = 'лосось'),
     (select id from saladRecipes.measurement where name = 'граммов'),
     9),

    ((select id from saladRecipes.recipe where numberofservings = 5),
     (select id from saladRecipes.ingredient where name = 'говядина'),
     (select id from saladRecipes.measurement where name = 'граммов'),
     10),
    ((select id from saladRecipes.recipe where numberofservings = 5),
     (select id from saladRecipes.ingredient where name = 'яблоко'),
     (select id from saladRecipes.measurement where name = 'граммов'),
     11),
    ((select id from saladRecipes.recipe where numberofservings = 5),
     (select id from saladRecipes.ingredient where name = 'морковь'),
     (select id from saladRecipes.measurement where name = 'граммов'),
     12),
    ((select id from saladRecipes.recipe where numberofservings = 5),
     (select id from saladRecipes.ingredient where name = 'лосось'),
     (select id from saladRecipes.measurement where name = 'граммов'),
     13),
    ((select id from saladRecipes.recipe where numberofservings = 5),
     (select id from saladRecipes.ingredient where name = 'молоко'),
     (select id from saladRecipes.measurement where name = 'граммов'),
     14);

insert into saladRecipes.recipeStep(id, name, description, recipeid, stepnum)
values ('01000000-0000-0000-0000-000000000000', 'step', 'description', '02000000-0000-0000-0000-000000000000', 1),
       ('07000000-0000-0000-0000-000000000000', 'step', 'description', '03000000-0000-0000-0000-000000000000', 1),

       ('02000000-0000-0000-0000-000000000000', 'first', 'first', '01000000-0000-0000-0000-000000000000', 1),
       ('03000000-0000-0000-0000-000000000000', 'second', 'second', '01000000-0000-0000-0000-000000000000', 2),
       ('04000000-0000-0000-0000-000000000000', 'third', 'third', '01000000-0000-0000-0000-000000000000', 3),
       ('05000000-0000-0000-0000-000000000000', 'fourth', 'fourth', '01000000-0000-0000-0000-000000000000', 4),
       ('06000000-0000-0000-0000-000000000000', 'fifth', 'fifth', '01000000-0000-0000-0000-000000000000', 5),

       ('08000000-0000-0000-0000-000000000000', 'first', 'first', '04000000-0000-0000-0000-000000000000', 1),
       ('09000000-0000-0000-0000-000000000000', 'second', 'second', '04000000-0000-0000-0000-000000000000', 2),
       ('0a000000-0000-0000-0000-000000000000', 'third', 'third', '04000000-0000-0000-0000-000000000000', 3);

insert into saladRecipes.user(name, email, login, password)
values ('existingUser', 'existingMail@mail.ru', 'anotherUsername', 'pass');

