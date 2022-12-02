#include <filesystem>
#include <fstream>
#include <iostream>
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
    std::vector<int> max = {0, 0, 0};
    int curr = 0;
    for (std::string line; std::getline(input, line);) {
      if (line.empty()) {

        max.push_back(curr);
        std::sort(max.begin(), max.end(), std::greater<int>());
        max.resize(3);

        curr = 0;
        continue;
      }

      int lineVal;
      auto ok = absl::SimpleAtoi(line, &lineVal);
      if (!ok) {
        std::cerr << "bad line: " << line << std::endl;
        return 1;
      }

      curr += lineVal;
    }
    max.push_back(curr);
    std::sort(max.begin(), max.end(), std::greater<int>());
    max.resize(3);

    std::cout << "part 1: " << max[0] << std::endl;
    std::cout << "part 2: " << max[0] + max[1] + max[2] << std::endl;
  }

  return 0;
}
