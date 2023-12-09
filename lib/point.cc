#include "./point.h"

#include <iostream>
#include <vector>

#include "absl/strings/numbers.h"
#include "absl/strings/str_cat.h"
#include "absl/strings/str_split.h"

namespace aoc {

absl::StatusOr<Point> Point::FromString(absl::string_view str) {
  std::vector<absl::string_view> parts = absl::StrSplit(str, ',');
  if (parts.size() != 2) {
    return absl::InvalidArgumentError(absl::StrCat("Not a valid point: ", str));
  }

  int x, y;
  if (!absl::SimpleAtoi(parts[0], &x) || !absl::SimpleAtoi(parts[1], &y)) {
    return absl::InvalidArgumentError(
        absl::StrCat("couldn't parse ints: ", str));
  }
  return Point{.x = x, .y = y};
}

std::ostream& operator<<(std::ostream& os, const Point& p) {
  return os << "(" << p.x << ", " << p.y << ")";
}

}  // namespace aoc
