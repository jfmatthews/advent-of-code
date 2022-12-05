#include <filesystem>
#include <fstream>
#include <iostream>
#include <regex>
#include <vector>

#include "absl/strings/numbers.h"
#include "absl/strings/str_split.h"
#include "absl/strings/string_view.h"

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

    // load input data
    std::vector<std::vector<std::string>> stacks;
    for (std::string line; std::getline(input, line);) {
      if (line.empty()) break;
      stacks.push_back({});

      std::vector<absl::string_view> blocks = absl::StrSplit(line, " ");
      for (auto b : blocks) {
        stacks.back().push_back(std::string(b));
      }
    }

    // operate
    std::regex command_pattern("move ([0-9]+) from ([0-9]+) to ([0-9]+)");
    for (std::string line; std::getline(input, line);) {
      std::smatch matches;
      if (std::regex_match(line, matches, command_pattern) &&
          matches.size() == 4) {
        int count, from, to;
        if (!absl::SimpleAtoi(matches[1].str(), &count)) {
          std::cerr << "unparseable COUNT in command: " << line << std::endl;
        }
        if (!absl::SimpleAtoi(matches[2].str(), &from)) {
          std::cerr << "unparseable FROM in command: " << line << std::endl;
        }
        if (!absl::SimpleAtoi(matches[3].str(), &to)) {
          std::cerr << "unparseable TO in command: " << line << std::endl;
        }

        if (stacks[from - 1].size() < count) {
          std::cerr << "invalid command " << line << "; current size of "
                    << from << " is " << stacks[from - 1].size() << std::endl;
        }

        for (int i = 0; i < count; i++) {
          auto box = stacks[from - 1].back();
          stacks[from - 1].pop_back();
          stacks[to - 1].push_back(box);
        }

      } else {
        std::cerr << "unparseable command: " << line << std::endl;
      }
    }

    std::cout << "part 1: ";
    for (auto& stack : stacks) {
      std::cout << stack.back();
    }
    std::cout << std::endl;
  }
  {
    std::ifstream input(abs_inpath.string());

    // load input data
    std::vector<std::vector<std::string>> stacks;
    for (std::string line; std::getline(input, line);) {
      if (line.empty()) break;
      stacks.push_back({});

      std::vector<absl::string_view> blocks = absl::StrSplit(line, " ");
      for (auto b : blocks) {
        stacks.back().push_back(std::string(b));
      }
    }

    // operate
    std::regex command_pattern("move ([0-9]+) from ([0-9]+) to ([0-9]+)");
    for (std::string line; std::getline(input, line);) {
      std::smatch matches;
      if (std::regex_match(line, matches, command_pattern) &&
          matches.size() == 4) {
        int count, from, to;
        if (!absl::SimpleAtoi(matches[1].str(), &count)) {
          std::cerr << "unparseable COUNT in command: " << line << std::endl;
        }
        if (!absl::SimpleAtoi(matches[2].str(), &from)) {
          std::cerr << "unparseable FROM in command: " << line << std::endl;
        }
        if (!absl::SimpleAtoi(matches[3].str(), &to)) {
          std::cerr << "unparseable TO in command: " << line << std::endl;
        }

        if (stacks[from - 1].size() < count) {
          std::cerr << "invalid command " << line << "; current size of "
                    << from << " is " << stacks[from - 1].size() << std::endl;
        }

        auto& fromStack = stacks[from - 1];
        auto& toStack = stacks[to - 1];
        for (int i = count; i > 0; i--) {
          auto box = fromStack[fromStack.size() - i];
          toStack.push_back(box);
        }
        fromStack.resize(fromStack.size() - count);
      } else {
        std::cerr << "unparseable command: " << line << std::endl;
      }
    }

    std::cout << "part 2: ";
    for (auto& stack : stacks) {
      std::cout << stack.back();
    }
    std::cout << std::endl;
  }
  return 0;
}
