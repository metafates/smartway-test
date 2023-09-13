insert into provider (id, name) values
    ('AA', 'AmericanAir'),
    ('IF', 'InternationFlights'),
    ('RS', 'RedStar');

insert into airline (code, name) values
    ('SU', 'Аэрофлот'),
    ('S7', 'S7'),
    ('KV', 'КрасАвиа'),
    ('U6', 'Уральские авиалинии'),
    ('UT', 'ЮТэйр'),
    ('FZ', 'Flydubai'),
    ('JB', 'JetBlue'),
    ('SJ', 'SuperJet'),
    ('WZ', 'Wizz Air'),
    ('N4', 'Nordwind Airlines'),
    ('5N', 'SmartAvia');

insert into schema (id, name) values
    (1, 'Primary'),
    (2, 'Test');

insert into account (id) values
    (1), -- demo, test schema
    (2), -- dev, test schema
    (3); -- sales, primary schema

insert into schema_provider (schema_id, provider_id) values
    (1, 'AA'),
    (1, 'IF'),
    (1, 'RS'),

    (2, 'IF'),
    (2, 'RS');

insert into account_schema (account_id, schema_id) values
    (1, 2),
    (2, 2),
    (3, 1);

insert into provider_airline (provider_id, airline_code) values
    ('AA', 'FZ'),
    ('AA', 'JB'),
    ('AA', 'SJ'),

    ('IF', 'SU'),
    ('IF', 'S7'),
    ('IF', 'FZ'),
    ('IF', 'N4'),
    ('IF', 'JB'),
    ('IF', 'WZ'),

    ('RS', 'SU'),
    ('RS', 'S7'),
    ('RS', 'KV'),
    ('RS', 'U6'),
    ('RS', 'UT'),
    ('RS', 'N4'),
    ('RS', '5N');
