#include <filesystem>
#include <fstream>
#include <iostream>
#include <memory>
#include <unordered_set>
#include <vector>

#include "absl/memory/memory.h"
#include "absl/strings/numbers.h"
#include "absl/strings/str_split.h"
#include "absl/strings/string_view.h"

class Range {
 public:
  static std::unique_ptr<Range> FromString(absl::string_view rep);

  bool Includes(int val) const { return left_ <= val && right_ >= val; }

  bool Subsumes(const Range& other) const {
    return left_ <= other.left_ && right_ >= other.right_;
  }

  bool Intersects(const Range& other) const {
    return Includes(other.left_) || Includes(other.right_) ||
           other.Includes(left_) || other.Includes(right_);
  }

  friend std::ostream& operator<<(std::ostream&, const Range&);

 private:
  Range(int l, int r) : left_(l), right_(r) {}

  // inclusive
  int left_;
  int right_;
};
std::ostream& operator<<(std::ostream& os, const Range& r) {
  os << "[" << r.left_ << ", " << r.right_ << "]";
  return os;
}

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
    int total = 0;

    for (std::string l; std::getline(input, l);) {
      std::pair<absl::string_view, absl::string_view> ranges =
          absl::StrSplit(l, ",");
      auto a = Range::FromString(ranges.first);
      auto b = Range::FromString(ranges.second);
      if (a == nullptr || b == nullptr) {
        std::cerr << "couldn't parse line " << l << std::endl;
        return 1;
      }

      if (a->Subsumes(*b) || b->Subsumes(*a)) {
        //  std::cout << "found total overlap on " << l << " (parsed as ranges "
        //     << *a << " and " << *b << ")" << std::endl;
        total++;
      }
    }

    std::cout << "part 1: " << total << std::endl;
  }
  {
    std::ifstream input(abs_inpath.string());
    int total = 0;

    for (std::string l; std::getline(input, l);) {
      std::pair<absl::string_view, absl::string_view> ranges =
          absl::StrSplit(l, ",");
      auto a = Range::FromString(ranges.first);
      auto b = Range::FromString(ranges.second);
      if (a == nullptr || b == nullptr) {
        std::cerr << "couldn't parse line " << l << std::endl;
        return 1;
      }

      if (a->Intersects(*b)) {
        //        std::cout << "found partial overlap on " << l << " (parsed as
        //        ranges "
        //                  << *a << " and " << *b << ")" << std::endl;
        total++;
      }
    }

    std::cout << "part 2: " << total << std::endl;
  }
}

std::unique_ptr<Range> Range::FromString(absl::string_view rep) {
  std::pair<absl::string_view, absl::string_view> limits =
      absl::StrSplit(rep, "-");

  int left, right;
  bool ok = absl::SimpleAtoi(limits.first, &left);
  if (!ok) {
    return nullptr;
  }

  ok = absl::SimpleAtoi(limits.second, &right);
  if (!ok) {
    return nullptr;
  }

  return absl::WrapUnique(new Range(left, right));
}
