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

constexpr const bool debug = false;

template <typename... Ts>
void l(Ts... ts) {
  if (debug) {
    (std::cout << ... << ts) << std::endl;
  }
}

struct Point {
  int x;
  int y;

  friend auto operator<=>(const Point&, const Point&) = default;
  friend std::ostream& operator<<(std::ostream&, const Point&);
};
std::ostream& operator<<(std::ostream& os, const Point& p) {
  return os << "(" << p.x << ", " << p.y << ")";
}

class Knot {
 public:
  Knot(Point init, Knot* next) : location_(init), next_(next) {}

  void MoveTo(Point new_loc);
  Point location() const { return location_; }

 private:
  void Tension();
  Point location_;
  Knot* next_;
};

void Knot::MoveTo(Point new_loc) {
  location_ = new_loc;
  Tension();
}

void Knot::Tension() {
  if (next_ == nullptr) {
    return;
  }

  Point next_new_loc = next_->location();
  const int dx = location_.x - next_new_loc.x;
  const int dy = location_.y - next_new_loc.y;

  if (std::abs(dx) > 1 && dy == 0) {
    next_new_loc.x += (dx > 0 ? 1 : -1);
  } else if (std::abs(dy) > 1 && dx == 0) {
    next_new_loc.y += (dy > 0 ? 1 : -1);
  } else if (std::abs(dy) > 1 || std::abs(dx) > 1) {
    next_new_loc.x += (dx > 0 ? 1 : -1);
    next_new_loc.y += (dy > 0 ? 1 : -1);
  }

  if (next_new_loc != next_->location()) {
    next_->MoveTo(next_new_loc);
  }
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

  std::vector<std::unique_ptr<Knot>> part1(2);
  part1[0] = std::make_unique<Knot>(Point{0, 0}, nullptr);
  for (int i = 1; i < part1.size(); i++) {
    part1[i] = std::make_unique<Knot>(Point{0, 0}, part1[i - 1].get());
  }
  std::vector<std::unique_ptr<Knot>> part2(10);
  part2[0] = std::make_unique<Knot>(Point{0, 0}, nullptr);
  for (int i = 1; i < part2.size(); i++) {
    part2[i] = std::make_unique<Knot>(Point{0, 0}, part2[i - 1].get());
  }

  std::set<Point> visited_p1{Point{0, 0}};
  std::set<Point> visited_p2{Point{0, 0}};
  for (std::string line; std::getline(input, line);) {
    std::vector<absl::string_view> s = absl::StrSplit(line, " ");
    if (s.size() != 2) {
      std::cerr << "bad line: " << line << std::endl;
      return 1;
    }
    int d;
    bool ok = absl::SimpleAtoi(s[1], &d);
    if (!ok) {
      std::cerr << "bad line: " << line << std::endl;
      return 1;
    }

    Knot* head1 = part1.back().get();
    Knot* head2 = part2.back().get();
    for (int i = 0; i < d; i++) {
      if (s[0] == "R") {
        l("right");
        head1->MoveTo(Point{
            .x = head1->location().x + 1,
            .y = head1->location().y,
        });
        head2->MoveTo(Point{
            .x = head2->location().x + 1,
            .y = head2->location().y,
        });
      } else if (s[0] == "L") {
        l("left");
        head1->MoveTo(Point{
            .x = head1->location().x - 1,
            .y = head1->location().y,
        });
        head2->MoveTo(Point{
            .x = head2->location().x - 1,
            .y = head2->location().y,
        });
      } else if (s[0] == "U") {
        l("up");
        head1->MoveTo(Point{
            .x = head1->location().x,
            .y = head1->location().y - 1,
        });
        head2->MoveTo(Point{
            .x = head2->location().x,
            .y = head2->location().y - 1,
        });
      } else if (s[0] == "D") {
        l("down");
        head1->MoveTo(Point{
            .x = head1->location().x,
            .y = head1->location().y + 1,
        });
        head2->MoveTo(Point{
            .x = head2->location().x,
            .y = head2->location().y + 1,
        });
      } else {
        std::cerr << "bad line: " << line << std::endl;
        return 1;
      }

      if (visited_p1.insert(part1.front()->location()).second) {
        l("p1 visited ", part1.front()->location());
      }
      if (visited_p2.insert(part2.front()->location()).second) {
        l("p2 visited ", part2.front()->location());
      }
    }

    std::cout << "after " << line << std::endl;
    for (auto& k : part2) {
      std::cout << k->location() << ", ";
    }
    std::cout << std::endl;
  }

  std::cout << "part 1: " << visited_p1.size() << std::endl;
  std::cout << "part 2: " << visited_p2.size() << std::endl;
  return 0;
}
