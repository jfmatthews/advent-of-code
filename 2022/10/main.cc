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

constexpr const bool debug = true;

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

int main(int argc, char** argv) {
  if (argc != 2) {
    std::cerr << "ERROR: needs 1 arg" << std::endl;
    return 1;
  }

  const std::filesystem::path inpath(argv[1]);
  const auto abs_inpath = std::filesystem::absolute(inpath);

  std::cout << "reading from " << abs_inpath << std::endl;

  std::ifstream input(abs_inpath.string());

  int x = 1;
  int crt_cursor = 0;
  int cycles_done = 0;
  int next_cycle_threshold = 20;
  int score = 0;
  for (std::string line; std::getline(input, line);) {
    std::vector<absl::string_view> tokens = absl::StrSplit(line, " ");
    int new_val = x;
    int cycles_taken = 0;
    if (tokens[0] == "noop") {
      cycles_taken = 1;
    } else if (tokens[0] == "addx") {
      if (tokens.size() != 2) {
        return DIE("bad line, needed 2 tokens: ", line);
      }
      int arg;
      bool ok = absl::SimpleAtoi(tokens[1], &arg);
      if (!ok) {
        return DIE("bad line, can't parse argument: ", line);
      }

      cycles_taken = 2;
      new_val += arg;
    }

    if ((cycles_done + cycles_taken) >= next_cycle_threshold) {
      //      l("at ", next_cycle_threshold, " register holds ", x, " for a
      //      score of ",
      //        next_cycle_threshold * x);
      score += (next_cycle_threshold * x);
      next_cycle_threshold += 40;
    }

    for (int i = 0; i < cycles_taken; i++) {
      if (std::abs((crt_cursor % 40) - x) <= 1) {
        std::cout << "#";
      } else {
        std::cout << ".";
      }
      crt_cursor++;
      if (crt_cursor % 40 == 0) {
        std::cout << std::endl;
      }
    }

    cycles_done += cycles_taken;
    x = new_val;
  }
  std::cout << "part 1: " << score << std::endl;
  return 0;
}
