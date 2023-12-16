#include <filesystem>
#include <fstream>
#include <iostream>
#include <map>
#include <string_view>
#include <utility>
#include <vector>

#include "absl/base/casts.h"
#include "absl/container/flat_hash_set.h"
#include "absl/strings/numbers.h"
#include "absl/strings/str_join.h"
#include "absl/strings/str_split.h"
#include "lib/point.h"

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

    int lineNo = 1;
    std::vector<int64_t> cardCounts(207, 1);
    for (std::string line; std::getline(input, line);) {
      std::pair<std::string_view, std::string_view> parts =
          absl::StrSplit(line, ": ");
      auto [_, cards] = parts;

      parts = absl::StrSplit(cards, " | ");
      auto [winners, have] = parts;

      std::vector<absl::string_view> winnerStrs = absl::StrSplit(winners, " ");
      absl::flat_hash_set<int> winnerIds;
      for (auto winnerStr : winnerStrs) {
        if (winnerStr.length() == 0) {
          continue;
        }
        int winnerId;
        if (!absl::SimpleAtoi(winnerStr, &winnerId)) {
          std::cerr << "ERROR: couldn't parse " << winnerStr << std::endl;
          return 1;
        }
        winnerIds.insert(winnerId);
      }

      std::vector<absl::string_view> haveStrs = absl::StrSplit(have, " ");
      int matches = 0;
      for (auto haveStr : haveStrs) {
        if (haveStr.length() == 0) {
          continue;
        }
        int haveId;
        if (!absl::SimpleAtoi(haveStr, &haveId)) {
          std::cerr << "ERROR: couldn't parse " << haveStr << std::endl;
          return 1;
        }
        if (winnerIds.find(haveId) != winnerIds.end()) {
          matches++;
        }
      }

      int score = (matches > 0) ? 1 << (matches - 1) : 0;
      std::cout << "card " << lineNo << " has " << matches
                << " matches for a score of " << score << std::endl;
      part1Total += score;

      std::cout << absl::StrJoin(cardCounts, " ") << std::endl;

      std::cout << "we have " << cardCounts[lineNo - 1] << " of these cards"
                << std::endl;
      for (int i = 0; i < matches && lineNo + i < cardCounts.size(); i++) {
        cardCounts[lineNo + i] += cardCounts[lineNo - 1];
      }

      lineNo++;
    }

    std::cout << "part 1: " << part1Total << std::endl;
    for (int i = 0; i < 207; i++) {
      part2Total += cardCounts[i];
    }
    std::cout << "part 2: " << part2Total << std::endl;
  }

  return 0;
}
