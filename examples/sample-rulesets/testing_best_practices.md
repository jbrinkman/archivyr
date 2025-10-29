# Testing Best Practices

## Overview

This ruleset defines best practices for writing effective, maintainable tests across different testing levels.

## Testing Pyramid

### Structure Your Tests

- **Unit Tests (70%)**: Test individual functions and methods
- **Integration Tests (20%)**: Test component interactions
- **End-to-End Tests (10%)**: Test complete user workflows

### Focus on Unit Tests

- Unit tests are fast and provide quick feedback
- They're easier to maintain and debug
- They should form the foundation of your test suite

## Unit Testing

### Test One Thing at a Time

- Each test should verify a single behavior
- Keep tests focused and simple
- Use descriptive test names

```python
# Good
def test_calculate_discount_returns_discounted_price():
    result = calculate_discount(100, 0.2)
    assert result == 80

def test_calculate_discount_raises_error_for_invalid_rate():
    with pytest.raises(ValueError):
        calculate_discount(100, 1.5)

# Bad
def test_calculate_discount():
    # Testing multiple things
    assert calculate_discount(100, 0.2) == 80
    assert calculate_discount(100, 0) == 100
    with pytest.raises(ValueError):
        calculate_discount(100, 1.5)
```

### Use the AAA Pattern

- **Arrange**: Set up test data and conditions
- **Act**: Execute the code being tested
- **Assert**: Verify the results

```python
def test_user_creation():
    # Arrange
    username = "john_doe"
    email = "john@example.com"
    
    # Act
    user = User.create(username, email)
    
    # Assert
    assert user.username == username
    assert user.email == email
    assert user.is_active is True
```

### Test Edge Cases

- Test boundary conditions
- Test error conditions
- Test empty inputs
- Test null/None values

```python
def test_divide_by_zero_raises_error():
    with pytest.raises(ZeroDivisionError):
        divide(10, 0)

def test_process_empty_list_returns_empty_result():
    result = process_items([])
    assert result == []

def test_find_user_with_none_id_raises_error():
    with pytest.raises(ValueError):
        find_user(None)
```

### Use Test Fixtures

- Extract common setup into fixtures
- Keep tests DRY (Don't Repeat Yourself)
- Use fixtures for test data

```python
@pytest.fixture
def sample_user():
    return User(id=1, username="john_doe", email="john@example.com")

def test_user_can_update_email(sample_user):
    sample_user.update_email("newemail@example.com")
    assert sample_user.email == "newemail@example.com"
```

## Integration Testing

### Test Component Interactions

- Test how components work together
- Use real dependencies when possible
- Test database interactions
- Test API integrations

```python
def test_user_repository_saves_and_retrieves_user(db_session):
    # Arrange
    repo = UserRepository(db_session)
    user = User(username="john_doe", email="john@example.com")
    
    # Act
    repo.save(user)
    retrieved = repo.find_by_username("john_doe")
    
    # Assert
    assert retrieved.username == user.username
    assert retrieved.email == user.email
```

### Use Test Containers

- Use containers for external dependencies
- Ensure tests are isolated and repeatable
- Clean up after tests

```python
@pytest.fixture(scope="module")
def postgres_container():
    with PostgresContainer("postgres:15") as postgres:
        yield postgres

def test_database_operations(postgres_container):
    connection = create_connection(postgres_container.get_connection_url())
    # Test database operations
```

## End-to-End Testing

### Test User Workflows

- Test complete user journeys
- Use realistic test data
- Test critical paths first

```python
def test_user_registration_and_login_flow():
    # Register new user
    response = client.post("/register", json={
        "username": "john_doe",
        "email": "john@example.com",
        "password": "secure_password"
    })
    assert response.status_code == 201
    
    # Login with credentials
    response = client.post("/login", json={
        "username": "john_doe",
        "password": "secure_password"
    })
    assert response.status_code == 200
    assert "token" in response.json()
```

### Keep E2E Tests Minimal

- E2E tests are slow and brittle
- Focus on critical user paths
- Don't duplicate unit test coverage

## Test Organization

### File Structure

- Place tests near the code they test
- Use consistent naming conventions
- Group related tests

```
src/
  user/
    user.py
    user_test.py
  order/
    order.py
    order_test.py
tests/
  integration/
    test_user_repository.py
  e2e/
    test_user_workflows.py
```

### Test Naming

- Use descriptive names that explain what is tested
- Include the expected behavior
- Make failures easy to understand

```python
# Good
def test_user_creation_with_valid_data_succeeds()
def test_user_creation_with_duplicate_email_raises_error()
def test_inactive_user_cannot_login()

# Bad
def test_user_1()
def test_create()
def test_error()
```

## Mocking and Stubbing

### Mock External Dependencies

- Mock external APIs
- Mock slow operations
- Mock non-deterministic behavior

```python
def test_send_notification_calls_email_service(mocker):
    # Arrange
    mock_email = mocker.patch('app.services.email.send')
    user = User(email="john@example.com")
    
    # Act
    send_notification(user, "Welcome!")
    
    # Assert
    mock_email.assert_called_once_with(
        to="john@example.com",
        subject="Welcome!"
    )
```

### Don't Over-Mock

- Prefer real objects when possible
- Only mock at boundaries
- Avoid mocking internal implementation details

## Test Data

### Use Realistic Data

- Use data that resembles production
- Test with various data types
- Include edge cases in test data

### Use Factories or Builders

- Create test data with factories
- Make test data creation reusable
- Use builders for complex objects

```python
class UserFactory:
    @staticmethod
    def create(username="john_doe", email="john@example.com", **kwargs):
        return User(username=username, email=email, **kwargs)

def test_user_with_custom_data():
    user = UserFactory.create(username="jane_doe", is_admin=True)
    assert user.is_admin is True
```

## Test Coverage

### Aim for High Coverage

- Target 80%+ code coverage
- Focus on critical paths
- Don't chase 100% coverage blindly

### Coverage Doesn't Equal Quality

- High coverage doesn't guarantee good tests
- Focus on testing behavior, not lines of code
- Review coverage reports to find gaps

## Performance

### Keep Tests Fast

- Unit tests should run in milliseconds
- Integration tests in seconds
- E2E tests in minutes

### Parallelize Tests

- Run tests in parallel when possible
- Ensure tests are independent
- Use test isolation

## Continuous Integration

### Run Tests Automatically

- Run tests on every commit
- Run tests on pull requests
- Block merges if tests fail

### Test Different Environments

- Test on multiple platforms
- Test with different dependency versions
- Test with different configurations

## Best Practices Summary

1. **Write tests first** (TDD) or alongside code
2. **Keep tests simple** and focused
3. **Test behavior**, not implementation
4. **Use descriptive names** for tests
5. **Maintain test code** like production code
6. **Run tests frequently** during development
7. **Fix failing tests** immediately
8. **Review test coverage** regularly
9. **Refactor tests** when needed
10. **Document complex test scenarios**
