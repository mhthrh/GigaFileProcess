create user mohsen identified by mohsen;
grant dba,connect to mohsen;
create table mohsen.transactions(
    id FLOAT,
    fullname varchar2(100),
    s_iban varchar2(30),
    d_iban varchar2(30),
    amount FLOAT,
    datetime varchar2(100)
    )