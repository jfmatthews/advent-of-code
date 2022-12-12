#include <compare>
#include <filesystem>
#include <fstream>
#include <iostream>
#include <regex>
#include <vector>

#include "absl/memory/memory.h"
#include "absl/strings/numbers.h"
#include "absl/strings/str_cat.h"
#include "absl/strings/str_split.h"
#include "absl/strings/string_view.h"

constexpr const bool debug = false;

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

int main(int argc, char** argv) {
  if (argc != 2) {
    std::cerr << "ERROR: needs 1 arg" << std::endl;
    return 1;
  }

  const std::filesystem::path inpath(argv[1]);
  const auto abs_inpath = std::filesystem::absolute(inpath);

  std::cout << "reading from " << abs_inpath << std::endl;

  std::ifstream input(abs_inpath.string());

  // load map
  Point start;
  Point end;
  std::vector<std::vector<int>> heightmap;
  for (std::string line; std::getline(input, line);) {
    std::vector<int> row(line.length());
    for (int i = 0; i < line.length(); i++) {
      if (line[i] == 'S') {
        start = Point{.x = i, .y = (int)heightmap.size()};
        l("starting at ", start);
        row[i] = 0;
      } else if (line[i] == 'E') {
        end = Point{.x = i, .y = (int)heightmap.size()};
        l("ending at ", end);
        row[i] = 25;
      } else {
        row[i] = line[i] - 'a';
      }
    }
    heightmap.emplace_back(std::move(row));
  }
  l("done loading ", heightmap.size(), " lines");

  // initialize distances
  std::vector<std::vector<int>> d_from_end;
  for (auto& row : heightmap) {
    d_from_end.push_back(
        std::vector<int>(row.size(), std::numeric_limits<int>::max()));
  }

  // BFS
  l("starting search...");
  std::deque<std::pair<Point, int>> queue{{end, 0}};
  while (!queue.empty()) {
    auto next = queue.front();
    queue.pop_front();

    const auto& p = next.first;
    const int d = next.second;
    l("flooding from", next.first, " at distance ", next.second);
    if (d_from_end[p.y][p.x] <= d) {
      // already visited
      continue;
    }

    d_from_end[p.y][p.x] = d;
    for (const auto& n : p.neighbors()) {
      l("considering ", n, " as neighbor of ", p);
      if (n.x < 0 || n.x >= heightmap[0].size() || n.y < 0 ||
          n.y >= heightmap.size()) {
        // off the map
        continue;
      }
      if (heightmap[p.y][p.x] - heightmap[n.y][n.x] > 1) {
        continue;
      }
      l("traveling to ", n, " from ", p, " at distance ", d + 1);
      queue.push_back({n, d + 1});
    }
  }

  std::cout << "part 1: " << d_from_end[start.y][start.x] << std::endl;

  int min_distance = std::numeric_limits<int>::max();
  for (int y = 0; y < heightmap.size(); y++) {
    for (int x = 0; x < heightmap[0].size(); x++) {
      if (d_from_end[y][x] < min_distance && heightmap[y][x] == 0) {
        l("new best: ", d_from_end[y][x], " starting from (", x, ", ", y, ")");
        min_distance = d_from_end[y][x];
      }
    }
  }
  std::cout << "part 2: " << min_distance << std::endl;
}
