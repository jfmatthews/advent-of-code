#include <filesystem>
#include <fstream>
#include <iostream>
#include <vector>

#include "absl/strings/str_split.h"

int main(int argc, char **argv) {
  if (argc != 2) {
    std::cerr << "ERROR: needs 1 arg" << std::endl;
    return 1;
  }

  const std::filesystem::path inpath(argv[1]);
  const auto abs_inpath = std::filesystem::absolute(inpath);

  std::cout << "reading from " << abs_inpath << std::endl;

  const std::map<absl::string_view, absl::string_view> wins = {
      {/*rock*/ "A", /*scissors*/ "Z"},  //
      {/*paper*/ "B", /*rock*/ "X"},     //
      {/*scissors*/ "C", /*paper*/ "Y"}, //
  };
  const std::map<absl::string_view, absl::string_view> loses = {
      {/*rock*/ "A", /*paper*/ "Y"},     //
      {/*paper*/ "B", /*scissors*/ "Z"}, //
      {/*scissors*/ "C", /*rock*/ "X"},  //
  };
  const std::map<absl::string_view, absl::string_view> ties = {
      {/*rock*/ "A", /*rock*/ "X"},         //
      {/*paper*/ "B", /*paper*/ "Y"},       //
      {/*scissors*/ "C", /*scissors*/ "Z"}, //
  };
  const std::map<absl::string_view, int> pts = {
      {"X", 1}, //
      {"Y", 2}, //
      {"Z", 3}, //
  };

  {
    std::ifstream input(abs_inpath.string());
    int total = 0;
    for (std::string line; std::getline(input, line);) {
      std::vector<absl::string_view> chars = absl::StrSplit(line, " ");
      if (chars.size() != 2) {
        std::cerr << "bad line: " << line << std::endl;
      }

      int thisRound = pts.at(chars[1]);

      if (wins.at(chars[0]) == chars[1]) {
        // lose
      } else if (loses.at(chars[0]) == chars[1]) {
        // win
        thisRound += 6;
      } else {
        // tie
        thisRound += 3;
      }
      total += thisRound;
    }

    std::cout << "part 1 answer: " << total << std::endl;
  }

  {
    std::ifstream input(abs_inpath.string());
    int total = 0;
    for (std::string line; std::getline(input, line);) {
      std::vector<absl::string_view> chars = absl::StrSplit(line, " ");
      if (chars.size() != 2) {
        std::cerr << "bad line: " << line << std::endl;
      }

      absl::string_view actual_play;
      int thisRound = 0;
      if (chars[1] == "X") {
        actual_play = wins.at(chars[0]);
      } else if (chars[1] == "Y") {
        thisRound += 3;
        actual_play = ties.at(chars[0]);
      } else if (chars[1] == "Z") {
        thisRound += 6;
        actual_play = loses.at(chars[0]);
      }
      thisRound += pts.at(actual_play);

      total += thisRound;
    }
    std::cout << "part 2 answer: " << total << std::endl;
  }
  return 0;
}
