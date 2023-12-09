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

const absl::flat_hash_map<std::string, int> colorMaxes = {
    {"red", 12},
    {"green", 13},
    {"blue", 14},
};

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
    for (std::string line; std::getline(input, line);) {
      auto [gameIdStr, gamesStr] =
          absl::implicit_cast<std::pair<std::string, std::string>>(
              absl::StrSplit(line, ":"));
      std::pair<std::string, std::string> gameIdParts =
          absl::StrSplit(gameIdStr, " ");

      int gameIdNum;
      if (!absl::SimpleAtoi(gameIdParts.second, &gameIdNum)) {
        std::cerr << "ERROR: couldn't parse gameIdNum" << std::endl;
        return 1;
      }

      bool gamePossible = true;
      std::vector<absl::string_view> pulls = absl::StrSplit(gamesStr, "; ");
      absl::flat_hash_map<std::string, int> max_by_color;
      for (auto pull : pulls) {
        pull = absl::StripAsciiWhitespace(pull);
        std::vector<absl::string_view> colors = absl::StrSplit(pull, ", ");
        for (auto colorPull : colors) {
          auto [countStr, color] =
              absl::implicit_cast<std::pair<std::string, std::string>>(
                  absl::StrSplit(colorPull, " "));
          int count;
          if (!absl::SimpleAtoi(countStr, &count)) {
            std::cerr << "ERROR: couldn't parse count in game " << gameIdNum
                      << " pull " << colorPull << std::endl;
            return 1;
          }

          // for part 1
          auto itr = colorMaxes.find(color);
          if (itr != colorMaxes.end()) {
            if (count > itr->second) {
              gamePossible = false;
            }
          }

          // for part 2
          max_by_color[color] = std::max(max_by_color[color], count);
        }
      }
      if (gamePossible) {
        part1Total += gameIdNum;
      }

      // for part 2
      int power = 1;
      for (auto [color, max] : max_by_color) {
        power *= max;
      }
      part2Total += power;
    }

    std::cout << "part 1: " << part1Total << std::endl;
    std::cout << "part 2: " << part2Total << std::endl;
  }

  return 0;
}
