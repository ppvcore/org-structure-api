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
Content-Type: application/json

{
  "name": "Company"
}

#### 201 Created
{
  "id": 1,
  "name": "Company",
  "parent_id": null,
  "created_at": "2026-03-08T21:43:53.611581724Z"
}

### 2. Создать департамент внутри Company
POST http://localhost:8080/departments
Content-Type: application/json

{
  "name": "IT",
  "parent_id": 1
}

#### 201 Created
{
  "id": 2,
  "name": "IT",
  "parent_id": 1,
  "created_at": "2026-03-08T21:44:21.68173677Z"
}

### 3. Создать ещё один в IT (ID=2)
POST http://localhost:8080/departments
Content-Type: application/json

{
  "name": "Backend",
  "parent_id": 2
}

#### 201 Created
{
  "id": 3,
  "name": "Backend",
  "parent_id": 2,
  "created_at": "2026-03-08T21:44:39.624861073Z"
}

### 4. Попытка дубликата имени в том же родителе → должна быть ошибка
POST http://localhost:8080/departments
Content-Type: application/json

{
  "name": "Backend",
  "parent_id": 2
}

#### 400 Bad Request
department name must be unique in parent

### 5. Создать сотрудника в Backend (предполагаем ID=3)
POST http://localhost:8080/departments/3/employees
Content-Type: application/json

{
  "full_name": "Ivan Ivanov",
  "position": "Senior Go Developer",
  "hired_at": "2026-03-09"
}

#### 201 Created
{
  "id": 2,
  "department_id": 3,
  "full_name": "Ivan Ivanov",
  "position": "Senior Go Developer",
  "hired_at": "2026-03-09T00:00:00Z",
  "created_at": "2026-03-08T21:46:03.316408259Z"
}

### 6. Получить дерево с глубиной 3 и сотрудниками
GET http://localhost:8080/departments/1?depth=3&include_employees=true

#### 200 OK
{
  "id": 1,
  "name": "Company",
  "parent_id": null,
  "created_at": "2026-03-08T21:43:53.611581Z",
  "children": [
    {
      "id": 2,
      "name": "IT",
      "parent_id": 1,
      "created_at": "2026-03-08T21:44:21.681736Z",
      "children": [
        {
          "id": 3,
          "name": "Backend",
          "parent_id": 2,
          "created_at": "2026-03-08T21:44:39.624861Z",
          "employees": [
            {
              "id": 2,
              "department_id": 3,
              "full_name": "Ivan Ivanov",
              "position": "Senior Go Developer",
              "hired_at": "2026-03-09T00:00:00Z",
              "created_at": "2026-03-08T21:46:03.316408Z"
            }
          ]
        }
      ]
    }
  ]
}

### 7. Переименовать + переместить Backend в корень
PATCH http://localhost:8080/departments/3
Content-Type: application/json

{
  "name": "Core Backend",
  "parent_id": 1
}

#### 200 OK
{
  "id": 3,
  "name": "Core Backend",
  "parent_id": 1,
  "created_at": "2026-03-08T21:44:39.624861Z"
}

### 8. Попытка цикла (переместить Company внутрь Backend) → ошибка
PATCH http://localhost:8080/departments/1
Content-Type: application/json

{
  "parent_id": 3
}

#### 400 Bad Request
would create cycle in department tree

### 9. Удалить IT с переносом сотрудников (предполагаем новый отдел ID=4)
DELETE http://localhost:8080/departments/2?mode=reassign&reassign_to_department_id=4

#### 204 No Content

### 10. Каскадное удаление тестового департамента
DELETE http://localhost:8080/departments/3?mode=cascade

#### 200 No Content