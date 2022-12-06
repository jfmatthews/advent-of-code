#include <filesystem>
#include <fstream>
#include <iostream>
#include <regex>
#include <vector>

#include "absl/strings/numbers.h"
#include "absl/strings/str_split.h"
#include "absl/strings/string_view.h"

int findHeader(absl::string_view, int);

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
    std::string line;
    std::getline(input, line);

    std::cout << "part 1: " << findHeader(line, 4) << std::endl;
    std::cout << "part 2: " << findHeader(line, 14) << std::endl;
  }

  return 0;
}

int findHeader(absl::string_view str, int length) {
  std::vector<int> last_seen(26, -1);
  int currStart = 0;
  for (int i = 0; i < str.length(); i++) {
    char c = str[i];
    if (last_seen[c - 'a'] >= currStart) {
      currStart = last_seen[c - 'a'] + 1;
    }
    last_seen[c - 'a'] = i;
    if (i - currStart + 1 >= length) {
      // answers are 1-indexed
      return i + 1;
    }
  }
  return -1;
}
