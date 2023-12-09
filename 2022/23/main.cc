#include <filesystem>
#include <fstream>
#include <iostream>
#include <map>
#include <vector>

#include "lib/logging.h"
#include "lib/point.h"

namespace aoc {

struct Elf {
  bool moving = false;
  Point proposed;
};

void RunPart1(std::ifstream& input) {
  std::map<Point, Elf> elfs;

  {
    int y = 0;
    for (std::string line; std::getline(input, line); y++) {
      for (int x = 0; x < line.length(); x++) {
        if (line[x] == '#') {
          Point loc{.x = x, .y = y};
          elfs[loc] = Elf{};
        }
      }
    }
  }

  std::cout << "Handling " << elfs.size() << " elfs." << std::endl;

  for (int step = 1; step <= 10; step++) {
    l("====ROUND ", step, "====");
    std::map<Point, int> proposals;
    for (auto& [p, elf] : elfs) {
      bool north_clear = true;
      bool south_clear = true;
      bool east_clear = true;
      bool west_clear = true;
      for (int dx : {-1, 0, 1}) {
        if (elfs.contains(Point{.x = p.x + dx, .y = p.y - 1})) {
          north_clear = false;
          break;
        }
      }
      for (int dx : {-1, 0, 1}) {
        if (elfs.contains(Point{.x = p.x + dx, .y = p.y + 1})) {
          south_clear = false;
          break;
        }
      }
      for (int dy : {-1, 0, 1}) {
        if (elfs.contains(Point{.x = p.x + 1, .y = p.y + dy})) {
          east_clear = false;
          break;
        }
      }
      for (int dy : {-1, 0, 1}) {
        if (elfs.contains(Point{.x = p.x - 1, .y = p.y + dy})) {
          west_clear = false;
          break;
        }
      }
      if (north_clear && south_clear && west_clear && east_clear) {
        elf.moving = false;
      } else if (north_clear) {
        elf.moving = true;
        auto to = Point{.x = p.x, .y = p.y - 1};
        elf.proposed = to;
        proposals[to]++;
      } else if (south_clear) {
        elf.moving = true;
        auto to = Point{.x = p.x, .y = p.y + 1};
        elf.proposed = to;
        proposals[to]++;
      } else if (east_clear) {
        elf.moving = true;
        auto to = Point{.x = p.x + 1, .y = p.y};
        elf.proposed = to;
        proposals[to]++;
      } else if (west_clear) {
        elf.moving = true;
        auto to = Point{.x = p.x - 1, .y = p.y};
        elf.proposed = to;
        proposals[to]++;
      } else {
        elf.moving = false;
      }
    }

    std::map<Point, Elf> next_elfs;
    for (auto& [p, elf] : elfs) {
      if (!elf.moving) {
        next_elfs[p] = elf;
      } else if (proposals[elf.proposed] == 1) {
        l("elf moving from ", p, " to ", elf.proposed);
        next_elfs[elf.proposed] = elf;
      } else {
        l("collision on ", elf.proposed, "; staying at ", p);
        next_elfs[p] = elf;
      }
    }
    next_elfs.swap(elfs);
  }

  int min_x = std::numeric_limits<int>::max();
  int min_y = std::numeric_limits<int>::max();
  int max_x = std::numeric_limits<int>::min();
  int max_y = std::numeric_limits<int>::min();
  for (auto& [p, _] : elfs) {
    min_x = std::min(min_x, p.x);
    max_x = std::max(max_x, p.x);
    min_y = std::min(min_y, p.y);
    max_y = std::max(max_y, p.y);
  }

  int box_a = (max_x - min_x) * (max_y - min_y);
  int empty_spaces = box_a - elfs.size();

  std::cout << "part 1: " << empty_spaces << std::endl;
}

void RunPart2(std::ifstream& input) {
  std::cout << "part 2: "
            << "..." << std::endl;
}
}  // namespace aoc

int main(int argc, char** argv) {
  if (argc != 2) {
    std::cerr << "ERROR: needs 1 arg" << std::endl;
    return 1;
  }

  const std::filesystem::path inpath(argv[1]);
  const auto abs_inpath = std::filesystem::absolute(inpath);

  std::cout << "reading from " << abs_inpath << std::endl;

  {
    std::ifstream input(abs_inpath.string());
    aoc::RunPart1(input);
  }

  {
    std::ifstream input(abs_inpath.string());
    aoc::RunPart2(input);
  }

  return 0;
}
