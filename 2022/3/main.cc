#include <filesystem>
#include <fstream>
#include <iostream>
#include <unordered_set>
#include <vector>

#include "absl/strings/str_split.h"
#include "absl/strings/string_view.h"

int score(char c) {
  if (c >= 'a' && c <= 'z') {
    return c - 'a' + 1;
  }
  if (c >= 'A' && c <= 'Z') {
    return c - 'A' + 27;
  }
  // unreachable with reasonable input
  std::cerr << "got bad char " << c << std::endl;
  return -1;
}

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
    int total = 0;

    for (std::string l; std::getline(input, l);) {
      absl::string_view line(l);
      auto left = std::string(line.substr(0, line.size() / 2));
      auto right = std::string(line.substr(line.size() / 2));
      //      std::cout << "left (" << left.length() << "): " << left << ",right
      //      ("
      //                << right.length() << "): " << right << std::endl;

      int thisRow = 0;
      int hits = 0;
      std::sort(left.begin(), left.end());
      std::sort(right.begin(), right.end());
      //      std::cout << "SORTED: left (" << left.length() << "): " << left
      //                << ",right (" << right.length() << "): " << right <<
      //                std::endl;
      int i = 0, j = 0;
      while (i < left.size() && j < right.size()) {
        if (left[i] == right[j]) {
          //          std::cout << "hit on " << left[i] << std::endl;
          char c = left[i];
          thisRow += score(c);
          hits++;
          while (left[i] == c) {
            i++;
          }
          while (right[j] == c) {
            j++;
          }
        }

        else if (left[i] < right[j]) {
          i++;
        } else if (left[i] > right[j]) {
          j++;
        }
      }

      std::cout << hits << " hits for total score of " << thisRow << std::endl;
      total += thisRow;
    }

    std::cout << "part 1: " << total << std::endl;
  }
  {
    std::ifstream input(abs_inpath.string());
    int total = 0;

    std::unordered_set<char> options;
    int rowsInSet = 0;
    for (std::string l; std::getline(input, l); rowsInSet++) {
      rowsInSet = rowsInSet % 3;

      std::unordered_set<char> thisRowOptions;
      for (char c : l) {
        thisRowOptions.insert(c);
      }

      if (rowsInSet == 0) {
        options.swap(thisRowOptions);
        continue;
      }

      std::unordered_set<char> new_options = options;
      for (char opt : options) {
        if (!thisRowOptions.contains(opt)) {
          new_options.erase(opt);
        }
      }
      options.swap(new_options);

      if (rowsInSet == 1) {
        continue;
      }

      if (options.size() != 1) {
        std::cout << "on line " << l
                  << ": should be down to 1 option, but have " << options.size()
                  << std::endl;
        return 1;
      }

      auto thisChunk = score(*options.begin());
      std::cout << "shared char is " << *options.begin() << " for score of "
                << thisChunk << std::endl;

      total += thisChunk;
      options.clear();
    }

    std::cout << "part 2: " << total << std::endl;
  }

  return 0;
}
