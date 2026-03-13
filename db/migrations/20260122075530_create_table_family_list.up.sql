create table family_list (
                             fl_id int primary key auto_increment,
                             cst_id int not null,
                             fl_relation varchar(50) not null,
                             fl_name varchar(50) not null,
                             fl_dob varchar(50) not null,
                             created_at timestamp default current_timestamp,
                             updated_at timestamp default current_timestamp on update current_timestamp,
                             foreign key (cst_id) references customers(cst_id)
) engine = InnoDB;