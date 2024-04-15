from faker import Faker
import random

faker = Faker('ru_RU')


def generate_users(num_users):
    users = []
    for i in range(num_users):
        user_full_name = faker.last_name() + ' ' + faker.first_name() + ' ' + faker.middle_name()
        email = faker.email()
        password = faker.password(length=10, special_chars=True, digits=True, upper_case=True, lower_case=True)
        users.append((user_full_name, email, password))
    return users


def generate_customers(num_customers, users, id):
    customers = []
    for i in range(num_customers):
        customer_id = id + i
        account_number = faker.random_number(digits=14)
        first_name = faker.first_name()
        last_name = faker.last_name()
        middle_name = faker.middle_name()
        date_of_birth = faker.date_of_birth(minimum_age=18, maximum_age=90).strftime('%Y-%m-%d')
        tax_id = faker.random_number(digits=10)
        user_full_name = random.choice(users)[0]
        status = 'Не в работе'
        customers.append(
            (customer_id, account_number, first_name, last_name, middle_name, date_of_birth, tax_id, user_id, status))
    return customers


def save_to_text_file(data, filename):
    with open(filename, 'w') as file:
        for row in data:
            if filename == 'customers.txt':
                row_str = f"('{row[0]}','{row[1]}','{row[2]}','{row[3]}','{row[4]}','{row[5]}','{row[6]}','{row[7]}', '{row[8]}')"
            else:
                row_str = f"('{row[0]}','{row[1]}','{row[2]}')"
            file.write(row_str + ',\n')


users = generate_users(1000)
customers = generate_customers(10000, users,1)

save_to_text_file(users, 'users.txt')
save_to_text_file(customers, 'customers.txt')
