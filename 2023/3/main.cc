#include <filesystem>
#include <fstream>
#include <iostream>
#include <map>
#include <utility>
#include <vector>

#include "absl/base/casts.h"
#include "absl/container/flat_hash_map.h"
#include "absl/strings/numbers.h"
#include "absl/strings/str_split.h"
#include "lib/point.h"

struct Component {
  int value;
  aoc::Point left;
  aoc::Point right;

  std::vector<aoc::Point> neighbors(int height, int width) const {
    std::vector<aoc::Point> candidates;
    candidates.reserve(((right.x - left.x) * 2) + 8);

    // up and down
    for (int x = left.x; x <= right.x; ++x) {
      candidates.push_back(aoc::Point{x, left.y - 1});
      candidates.push_back(aoc::Point{x, left.y + 1});
    }

    // left and right
    for (int y = left.y - 1; y <= left.y + 1; ++y) {
      candidates.push_back(aoc::Point{left.x - 1, y});
      candidates.push_back(aoc::Point{right.x + 1, y});
    }

    std::vector<aoc::Point> result;
    for (const auto &p : candidates) {
      if (p.x >= 0 && p.x < width && p.y >= 0 && p.y < height) {
        result.push_back(p);
      }
    }
    return result;
  }
};

bool isDigit(char c) { return c >= '0' && c <= '9'; }

bool isSymbol(char c) { return c != '.' && !isDigit(c); }

int main(int argc, char **argv) {
  if (argc != 2) {
    std::cerr << "ERROR: needs 1 arg" << std::endl;
    return 1;
  }

  const std::filesystem::path inpath(argv[1]);
  const auto abs_inpath = std::filesystem::absolute(inpath);

  std::cout << "reading from " << abs_inpath << std::endl;

  {
    std::ifstream input(abs_inpath.string());
    int part1Total = 0;
    int part2Total = 0;

    std::vector<Component> components;
    std::map<aoc::Point, char> symbols;

    Component currComponent;
    int lineNum = 0;
    int width = 0;
    for (std::string line; std::getline(input, line);) {
      width = line.size();
      bool in_component = false;
      for (int i = 0; i < line.size(); i++) {
        if (isDigit(line[i])) {
          if (!in_component) {
            currComponent = Component{
                .value = line[i] - '0',
                .left = aoc::Point{i, lineNum},
            };
            in_component = true;
          } else {
            currComponent.value *= 10;
            currComponent.value += line[i] - '0';
          }
        }

        // finish components
        if (in_component && !isDigit(line[i])) {
          currComponent.right = aoc::Point{i - 1, lineNum};
          components.push_back(currComponent);
          std::cout << "found component with val " << currComponent.value
                    << " ending at " << currComponent.right << std::endl;
          currComponent = Component{};
          in_component = false;
        }

        // blank
        if (line[i] == '.')
          continue;

        // symbols
        if (isSymbol(line[i])) {
          symbols[aoc::Point{i, lineNum}] = line[i];
          //           std::cout << "found symbol at " << aoc::Point{i, lineNum}
          //                     << std::endl;
        }
      }

      // finish components ending at end of line
      if (in_component) {
        currComponent.right =
            aoc::Point{static_cast<int>(line.size() - 1), lineNum};
        components.push_back(currComponent);
        //         std::cout << "found component with val " <<
        //         currComponent.value
        //                   << " ending at " << currComponent.right <<
        //                   std::endl;
        currComponent = Component{};
        in_component = false;
      }

      lineNum++;
    }

    for (const auto &c : components) {
      for (const auto &p : c.neighbors(lineNum, width)) {
        if (symbols.find(p) != symbols.end()) {
          std::cout << "component " << c.value << " at " << c.left
                    << " has neighbor at " << p << " with symbol "
                    << symbols.at(p) << std::endl;
          part1Total += c.value;
          break;
        }
      }
    }
    std::cout << "part 1: " << part1Total << std::endl;

    std::map<aoc::Point, std::vector<Component>> componentsByGear;
    for (const auto &c : components) {
      for (const auto &p : c.neighbors(lineNum, width)) {
        if (symbols.find(p) != symbols.end()) {
          if (symbols.at(p) == '*') {
            componentsByGear[p].push_back(c);
          }
        }
      }
    }

    for (const auto &[p, cs] : componentsByGear) {
      if (cs.size() == 2) {
        part2Total += (cs[0].value * cs[1].value);
      }
    }
    std::cout << "part 2: " << part2Total << std::endl;
  }

  return 0;
}
