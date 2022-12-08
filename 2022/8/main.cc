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

struct Tree {
  int height;
  int max_left;
  int max_right;
  int max_up;
  int max_down;

  bool visible() {
    return height > max_up ||     //
           height > max_down ||   //
           height > max_right ||  //
           height > max_left      //
        ;
  }
};

int main(int argc, char** argv) {
  if (argc != 2) {
    std::cerr << "ERROR: needs 1 arg" << std::endl;
    return 1;
  }

  const std::filesystem::path inpath(argv[1]);
  const auto abs_inpath = std::filesystem::absolute(inpath);

  std::cout << "reading from " << abs_inpath << std::endl;

  std::ifstream input(abs_inpath.string());

  std::vector<std::vector<Tree>> forest;
  for (std::string line; std::getline(input, line);) {
    forest.push_back(std::vector<Tree>());
    forest.back().reserve(line.length());
    for (char c : line) {
      forest.back().push_back(Tree{
          .height = c - '0',
          .max_left = -1,
          .max_right = -1,
          .max_up = -1,
          .max_down = -1,
      });
    }
  }

  // sweep left
  for (auto& row : forest) {
    for (int c = 1; c < row.size(); c++) {
      row[c].max_left = std::max(row[c - 1].height, row[c - 1].max_left);
    }
  }

  // sweep right
  for (auto& row : forest) {
    for (int c = row.size() - 2; c >= 0; c--) {
      row[c].max_right = std::max(row[c + 1].height, row[c + 1].max_right);
    }
  }

  // sweep top-down
  for (int r = 1; r < forest.size(); r++) {
    auto& row = forest[r];
    for (int c = 0; c < row.size(); c++) {
      row[c].max_up =
          std::max(forest[r - 1][c].height, forest[r - 1][c].max_up);
    }
  }

  // sweep bottom-up
  for (int r = forest.size() - 2; r >= 0; r--) {
    auto& row = forest[r];
    for (int c = 0; c < row.size(); c++) {
      row[c].max_down =
          std::max(forest[r + 1][c].height, forest[r + 1][c].max_down);
    }
  }

  // check visibility
  int count = 0;
  for (int r = 0; r < forest.size(); r++) {
    int rowVisible = 0;
    auto& row = forest[r];
    for (int c = 0; c < row.size(); c++) {
      if (row[c].visible()) rowVisible++;
    }
    std::cout << rowVisible << " of " << row.size() << " visible in row " << r
              << std::endl;
    count += rowVisible;
  }

  std::cout << "part 1: " << count << std::endl;

  // count scenic scores. scores for all edge trees is 0, since there's one
  // direction in which they see nothing
  int max_score = 0;
  for (int r = 1; r < forest.size() - 1; r++) {
    auto& row = forest[r];
    for (int c = 1; c < row.size() - 1; c++) {
      auto& tree = row[c];

      int up_score, down_score, left_score, right_score;
      if (tree.height > tree.max_up) {
        up_score = r;
      } else {
        up_score = 0;
        for (int r2 = r - 1; r2 >= 0; r2--) {
          up_score++;
          if (forest[r2][c].height >= tree.height) break;
        }
      }

      if (tree.height > tree.max_down) {
        down_score = (forest.size() - r - 1);
      } else {
        down_score = 0;
        for (int r2 = r + 1; r2 < forest.size(); r2++) {
          down_score++;
          if (forest[r2][c].height >= tree.height) break;
        }
      }

      if (tree.height > tree.max_left) {
        left_score = c;
      } else {
        left_score = 0;
        for (int c2 = c - 1; c2 >= 0; c2--) {
          left_score++;
          if (forest[r][c2].height >= tree.height) break;
        }
      }

      if (tree.height > tree.max_right) {
        right_score = row.size() - c - 1;
      } else {
        right_score = 0;
        for (int c2 = c + 1; c2 < row.size(); c2++) {
          right_score++;
          if (forest[r][c2].height >= tree.height) break;
        }
      }

      const int this_score = up_score * down_score * left_score * right_score;
      if (this_score > max_score) {
        std::cout << "new best score at " << r << ", " << c << ": "
                  << this_score << std::endl;
        max_score = this_score;
      }
    }
  }

  std::cout << "part 2: " << max_score << std::endl;

  return 0;
}
