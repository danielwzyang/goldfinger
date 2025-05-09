  // If snake (head) hits food, add +1 to the snake size and create a new food
  if (pos.x == food.x && pos.y == food.y) {
    len++;  // Increase length instead of size
    newFood();
  } 