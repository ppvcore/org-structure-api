# Organization Structure API

REST API для управления деревом подразделений и их сотрудниками

## Технологии

- Go 1.26
- PostgreSQL 18
- GORM
- Goose
- Docker-Compose
- testify
- net/http

## Функциональность

- Создание / обновление / удаление подразделений (с проверкой уникальности имён в пределах родителя и отсутствия циклов)
- Каскадное удаление или перенос сотрудников в другой отдел
- Получение подразделения с сотрудниками и вложенным деревом (с параметром depth)
- Создание сотрудников в конкретном подразделении

## Тесты

### Departments

- Создать с пустым именем → ошибка
- Создать с дублирующимся именем у родителя → ошибка
- Обновить с родителем = себе → ошибка

### Employees

- Создать с пустым full_name или position → ошибка
- Успешное создание → проверка полей и департамента
- Перевести в другой департамент → обновлён department_id
- Удалить по департаменту → сотрудники удалены

Покрыто: валидация, уникальность, дерево, каскадные операции

## Запуск тестов

go test ./...

## Запуск проекта (Docker-Compose)

docker compose up

## Результаты тестов API запросов

### 1. Создать корневой департамент

POST http://localhost:8080/departments
<br>
Content-Type: application/json<br>

{<br>
"name": "Company"<br>
}<br>

#### 201 Created

{<br>
"id": 1,<br>
"name": "Company",<br>
"parent_id": null,<br>
"created_at": "2026-03-08T21:43:53.611581724Z"<br>
}<br>

### 2. Создать департамент внутри Company

POST http://localhost:8080/departments
<br>
Content-Type: application/json<br>

{<br>
"name": "IT",<br>
"parent_id": 1<br>
}<br>

#### 201 Created

{<br>
"id": 2,<br>
"name": "IT",<br>
"parent_id": 1,<br>
"created_at": "2026-03-08T21:44:21.68173677Z"<br>
}<br>

### 3. Создать ещё один в IT (ID=2)

POST http://localhost:8080/departments
<br>
Content-Type: application/json<br>

{<br>
"name": "Backend",<br>
"parent_id": 2<br>
}<br>

#### 201 Created

{<br>
"id": 3,<br>
"name": "Backend",<br>
"parent_id": 2,<br>
"created_at": "2026-03-08T21:44:39.624861073Z"<br>
}<br>

### 4. Попытка дубликата имени в том же родителе → ошибка

POST http://localhost:8080/departments
<br>
Content-Type: application/json<br>

{<br>
"name": "Backend",<br>
"parent_id": 2<br>
}<br>

#### 400 Bad Request

department name must be unique in parent<br>

### 5. Создать сотрудника в Backend (предполагаем ID=3)

POST http://localhost:8080/departments/3/employees
<br>
Content-Type: application/json<br>

{<br>
"full_name": "Ivan Ivanov",<br>
"position": "Senior Go Developer",<br>
"hired_at": "2026-03-09"<br>
}<br>

#### 201 Created

{<br>
"id": 2,<br>
"department_id": 3,<br>
"full_name": "Ivan Ivanov",<br>
"position": "Senior Go Developer",<br>
"hired_at": "2026-03-09T00:00:00Z",<br>
"created_at": "2026-03-08T21:46:03.316408259Z"<br>
}<br>

### 6. Получить дерево с глубиной 3 и сотрудниками

GET http://localhost:8080/departments/1?depth=3&include_employees=true
<br>

#### 200 OK

{<br>
"id": 1,<br>
"name": "Company",<br>
"parent_id": null,<br>
"created_at": "2026-03-08T21:43:53.611581Z",<br>
"children": [<br>
{<br>
"id": 2,<br>
"name": "IT",<br>
"parent_id": 1,<br>
"created_at": "2026-03-08T21:44:21.681736Z",<br>
"children": [<br>
{<br>
"id": 3,<br>
"name": "Backend",<br>
"parent_id": 2,<br>
"created_at": "2026-03-08T21:44:39.624861Z",<br>
"employees": [<br>
{<br>
"id": 2,<br>
"department_id": 3,<br>
"full_name": "Ivan Ivanov",<br>
"position": "Senior Go Developer",<br>
"hired_at": "2026-03-09T00:00:00Z",<br>
"created_at": "2026-03-08T21:46:03.316408Z"<br>
}<br>
]<br>
}<br>
]<br>
}<br>
]<br>
}<br>

### 7. Переименовать + переместить Backend в корень

PATCH http://localhost:8080/departments/3
<br>
Content-Type: application/json<br>

{<br>
"name": "Core Backend",<br>
"parent_id": 1<br>
}<br>

#### 200 OK

{<br>
"id": 3,<br>
"name": "Core Backend",<br>
"parent_id": 1,<br>
"created_at": "2026-03-08T21:44:39.624861Z"<br>
}<br>

### 8. Попытка цикла (переместить Company внутрь Backend) → ошибка

PATCH http://localhost:8080/departments/1
<br>
Content-Type: application/json<br>

{<br>
"parent_id": 3<br>
}<br>

#### 400 Bad Request

would create cycle in department tree<br>

### 9. Удалить IT с переносом сотрудников (предполагаем новый отдел ID=3)

DELETE http://localhost:8080/departments/2?mode=reassign&reassign_to_department_id=3
<br>

#### 204 No Content<br>
### 10. Каскадное удаление тестового департамента

DELETE http://localhost:8080/departments/3?mode=cascade
<br>

#### 200 No Content<br>