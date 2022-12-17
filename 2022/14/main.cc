#include <compare>
#include <filesystem>
#include <fstream>
#include <iostream>
#include <regex>
#include <vector>

#include "absl/status/statusor.h"
#include "absl/strings/numbers.h"
#include "absl/strings/str_cat.h"
#include "absl/strings/str_split.h"
#include "absl/strings/string_view.h"

constexpr const bool debug = true;

template <typename... Ts>
void l(Ts... ts) {
  if (debug) {
    (std::cout << ... << ts) << std::endl;
  }
}

template <typename... Ts>
int DIE(Ts... ts) {
  (std::cerr << ... << ts) << std::endl;
  return 1;
}

struct Point {
  int x;
  int y;

  static absl::StatusOr<Point> FromString(absl::string_view str) {
    std::vector<absl::string_view> parts = absl::StrSplit(str, ',');
    if (parts.size() != 2) {
      return absl::InvalidArgumentError(
          absl::StrCat("Not a valid point: ", str));
    }

    int x, y;
    if (!absl::SimpleAtoi(parts[0], &x) || !absl::SimpleAtoi(parts[1], &y)) {
      return absl::InvalidArgumentError(
          absl::StrCat("couldn't parse ints: ", str));
    }
    return Point{.x = x, .y = y};
  }

  std::vector<Point> neighbors() const {
    return {
        {.x = x - 1, .y = y},
        {.x = x + 1, .y = y},
        {.x = x, .y = y - 1},
        {.x = x, .y = y + 1},
    };
  }

  friend auto operator<=>(const Point&, const Point&) = default;
  friend std::ostream& operator<<(std::ostream&, const Point&);
};
std::ostream& operator<<(std::ostream& os, const Point& p) {
  return os << "(" << p.x << ", " << p.y << ")";
}

void PrintMap(const std::set<Point>& walls, const std::set<Point>& sand,
              int min_x, int max_x, int max_y, bool with_wall = false) {
  for (int y = 0; y <= max_y; y++) {
    for (int x = min_x; x <= max_x; x++) {
      Point p{.x = x, .y = y};
      if (with_wall && y == max_y) {
        std::cout << "=";
      } else if (walls.contains(p)) {
        std::cout << "#";
      } else if (sand.contains(p)) {
        std::cout << "o";
      } else {
        std::cout << ".";
      }
    }
    std::cout << std::endl;
  }
}

void DrawLine(Point p1, Point p2, std::set<Point>* map) {
  l("Line from ", p1, " to ", p2);
  if (p2.x < p1.x || p2.y < p1.y) {
    std::swap(p1, p2);
  }

  for (int y = p1.y; y <= p2.y; y++) {
    map->insert(Point{.x = p1.x, .y = y});
    l("   point on line: ", Point{.x = p1.x, .y = y});
  }
  for (int x = p1.x; x <= p2.x; x++) {
    map->insert(Point{.x = x, .y = p1.y});
    l("   point on line: ", Point{.x = x, .y = p1.y});
  }
}

int main(int argc, char** argv) {
  if (argc != 2) {
    std::cerr << "ERROR: needs 1 arg" << std::endl;
    return 1;
  }

  const std::filesystem::path inpath(argv[1]);
  const auto abs_inpath = std::filesystem::absolute(inpath);

  std::cout << "reading from " << abs_inpath << std::endl;

  std::ifstream input(abs_inpath.string());

  int max_y = -1;
  int min_x = std::numeric_limits<int>::max();
  int max_x = -1;
  std::set<Point> walls;
  for (std::string line; std::getline(input, line);) {
    std::vector<absl::string_view> point_strs = absl::StrSplit(line, " -> ");
    std::vector<Point> points;
    for (auto p_str : point_strs) {
      if (auto parsed = Point::FromString(p_str); parsed.ok()) {
        max_y = std::max(max_y, parsed->y);
        min_x = std::min(min_x, parsed->x);
        max_x = std::max(max_x, parsed->x);

        points.push_back(*parsed);
      } else {
        std::cerr << "bad point: " << p_str << " in line " << line << std::endl;
        return 1;
      }
    }

    for (int i = 0; i < points.size() - 1; i++) {
      DrawLine(points[i], points[i + 1], &walls);
    }
  }

  {
    std::set<Point> sand;
    auto empty = [&](Point p) {
      return !walls.contains(p) && !sand.contains(p);
    };
    while (true) {
      Point p = {.x = 500, .y = 0};

      while (p.y < max_y) {
        if (empty(Point{.x = p.x, .y = p.y + 1})) {
          p.y++;
        } else if (empty(Point{.x = p.x - 1, .y = p.y + 1})) {
          p.x--;
          p.y++;
        } else if (empty(Point{.x = p.x + 1, .y = p.y + 1})) {
          p.x++;
          p.y++;
        } else {
          break;
        }
      }
      if (p.y >= max_y) {
        break;
      } else {
        min_x = std::min(min_x, p.x);
        max_x = std::max(max_x, p.x);
        sand.insert(p);
      }
    }

    if (debug) {
      PrintMap(walls, sand, min_x, max_x, max_y);
    }

    std::cout << "part 1: " << sand.size() << std::endl;
  }

  {
    std::set<Point> sand;
    const int bottom_wall = max_y + 2;
    auto empty = [&](Point p) {
      return p.y < bottom_wall && !walls.contains(p) && !sand.contains(p);
    };
    while (true) {
      Point p = {.x = 500, .y = 0};

      // guaranteed to terminate because y is capped at bottom_wall
      while (true) {
        if (empty(Point{.x = p.x, .y = p.y + 1})) {
          p.y++;
        } else if (empty(Point{.x = p.x - 1, .y = p.y + 1})) {
          p.x--;
          p.y++;
        } else if (empty(Point{.x = p.x + 1, .y = p.y + 1})) {
          p.x++;
          p.y++;
        } else {
          break;
        }
      }

      sand.insert(p);
      min_x = std::min(min_x, p.x);
      max_x = std::max(max_x, p.x);
      if (p == Point{.x = 500, .y = 0}) {
        break;
      }
    }

    if (debug) {
      PrintMap(walls, sand, min_x, max_x, bottom_wall, true);
    }

    std::cout << "part 2: " << sand.size() << std::endl;
  }

  return 0;
}
