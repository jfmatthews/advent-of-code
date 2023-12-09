#include <filesystem>
#include <fstream>
#include <iostream>
#include <map>
#include <vector>

#include "absl/strings/numbers.h"

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
    for (std::string line; std::getline(input, line);) {
      int tens = 0, ones = 0;
      for (int i = 0; i < line.size(); ++i) {
        if (line[i] >= '0' && line[i] <= '9') {
          tens = line[i] - '0';
          break;
        }
      }
      for (int i = line.size() - 1; i >= 0; --i) {
        if (line[i] >= '0' && line[i] <= '9') {
          ones = line[i] - '0';
          break;
        }
      }
      total += (tens * 10) + ones;
    }

    std::cout << "part 1: " << total << std::endl;
  }

  {
    const std::map<std::string, int> words = {
        {"zero", 0}, {"one", 1}, {"two", 2},   {"three", 3}, {"four", 4},
        {"five", 5}, {"six", 6}, {"seven", 7}, {"eight", 8}, {"nine", 9},
    };
    std::ifstream input(abs_inpath.string());
    int total = 0;
    for (std::string line; std::getline(input, line);) {
      int tens, ones;
      for (int i = 0; i < line.size(); ++i) {
        if (line[i] >= '0' && line[i] <= '9') {
          tens = line[i] - '0';
          goto foundTens;
        }
        for (const auto &[word, value] : words) {
          if (line.substr(i, word.size()) == word) {
            tens = value;
            goto foundTens;
          }
        }
      }
    foundTens:
      for (int i = line.size() - 1; i >= 0; --i) {
        if (line[i] >= '0' && line[i] <= '9') {
          ones = line[i] - '0';
          goto foundOnes;
        }
        for (const auto &[word, value] : words) {
          if (line.substr(i, word.size()) == word) {
            ones = value;
            goto foundOnes;
          }
        }
      }
    foundOnes:
      total += (tens * 10) + ones;
    }
    std::cout << "part 2: " << total << std::endl;
  }

  return 0;
}
