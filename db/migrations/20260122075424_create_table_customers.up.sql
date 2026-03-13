create table customers (
                           cst_id int primary key auto_increment,
                           nationality_id int not null,
                           cst_name char(50) not null,
                           cst_dob date not null,
                           cst_phoneNum varchar(20) not null,
                           cst_email varchar(50) not null,
                           created_at timestamp default current_timestamp,
                           updated_at timestamp default current_timestamp on update current_timestamp,
                           foreign key (nationality_id) references nationality(nationality_id),
                           unique key unique_email (cst_email)
) engine = InnoDB;
