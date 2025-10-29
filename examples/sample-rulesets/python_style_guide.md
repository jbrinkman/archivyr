# Python Style Guide

## Overview

This ruleset defines Python coding standards for team projects, based on PEP 8 and team-specific conventions.

## Naming Conventions

### Variables and Functions

- Use `snake_case` for variables and functions
- Use descriptive names that convey purpose
- Avoid single-letter names except for loop counters

```python
# Good
user_count = 0
def calculate_total_price(items):
    pass

# Bad
uc = 0
def calcTotPrice(items):
    pass
```

### Classes

- Use `PascalCase` for class names
- Use singular nouns for class names

```python
# Good
class UserAccount:
    pass

# Bad
class user_account:
    pass
```

### Constants

- Use `UPPER_SNAKE_CASE` for constants
- Define constants at module level

```python
# Good
MAX_RETRY_COUNT = 3
DEFAULT_TIMEOUT = 30

# Bad
maxRetryCount = 3
```

## Code Organization

### Imports

- Group imports in the following order:
  1. Standard library imports
  2. Related third-party imports
  3. Local application imports
- Use absolute imports when possible
- One import per line

```python
# Good
import os
import sys

import requests
from django.conf import settings

from myapp.models import User

# Bad
import os, sys
from myapp.models import *
```

### Function Length

- Keep functions under 50 lines when possible
- Extract complex logic into helper functions
- Each function should do one thing well

### Module Length

- Keep modules under 500 lines
- Split large modules into smaller, focused modules

## Documentation

### Docstrings

- Use docstrings for all public modules, functions, classes, and methods
- Follow Google or NumPy docstring format
- Include parameter types and return types

```python
def calculate_discount(price: float, discount_rate: float) -> float:
    """Calculate the discounted price.
    
    Args:
        price: The original price
        discount_rate: The discount rate (0.0 to 1.0)
        
    Returns:
        The discounted price
        
    Raises:
        ValueError: If discount_rate is not between 0 and 1
    """
    if not 0 <= discount_rate <= 1:
        raise ValueError("Discount rate must be between 0 and 1")
    return price * (1 - discount_rate)
```

## Error Handling

- Use specific exception types
- Don't use bare `except:` clauses
- Include meaningful error messages

```python
# Good
try:
    result = risky_operation()
except ValueError as e:
    logger.error(f"Invalid value: {e}")
    raise

# Bad
try:
    result = risky_operation()
except:
    pass
```

## Type Hints

- Use type hints for function parameters and return values
- Use `Optional` for nullable values
- Use `Union` sparingly

```python
from typing import Optional, List

def find_user(user_id: int) -> Optional[User]:
    """Find a user by ID."""
    pass

def get_user_names(users: List[User]) -> List[str]:
    """Extract names from user list."""
    return [user.name for user in users]
```

## Testing

- Write tests for all public functions
- Use descriptive test names
- Follow the Arrange-Act-Assert pattern
- Aim for 80%+ code coverage

```python
def test_calculate_discount_with_valid_rate():
    # Arrange
    price = 100.0
    discount_rate = 0.2
    
    # Act
    result = calculate_discount(price, discount_rate)
    
    # Assert
    assert result == 80.0
```

## Code Formatting

- Use `black` for automatic code formatting
- Line length: 88 characters (black default)
- Use 4 spaces for indentation (never tabs)

## Linting

- Use `pylint` or `flake8` for linting
- Use `mypy` for type checking
- Address all linting warnings before committing
