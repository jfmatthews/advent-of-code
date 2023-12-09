#ifndef LIB_POINT_H_
#define LIB_POINT_H_

#include <compare>
#include <iostream>
#include <vector>

#include "absl/status/statusor.h"
#include "absl/strings/string_view.h"

namespace aoc {

struct Point {
  int x;
  int y;

  static absl::StatusOr<Point> FromString(absl::string_view str);

  std::vector<Point> neighbors4() const {
    return {
        {.x = x - 1, .y = y},
        {.x = x + 1, .y = y},
        {.x = x, .y = y - 1},
        {.x = x, .y = y + 1},
    };
  }

  std::vector<Point> neighbors8() const {
    return {
        {.x = x - 1, .y = y - 1},  //
        {.x = x - 1, .y = y},      //
        {.x = x - 1, .y = y + 1},  //
        {.x = x, .y = y - 1},      //
        {.x = x, .y = y + 1},      //
        {.x = x + 1, .y = y - 1},  //
        {.x = x + 1, .y = y},      //
        {.x = x + 1, .y = y + 1},
    };
  }

  friend auto operator<=>(const Point&, const Point&) = default;
  friend std::ostream& operator<<(std::ostream&, const Point&);
};

std::ostream& operator<<(std::ostream& os, const Point& p);

}  // namespace aoc

#endif
