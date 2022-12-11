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

class Monkey;

class MonkeyBarrel {
 public:
  Monkey* GetMonkey(int idx) const {
    if (idx >= monkeys_.size()) {
      return nullptr;
    }
    return monkeys_[idx];
  }

  void AddMonkey(Monkey* m) { monkeys_.push_back(m); }

 private:
  std::vector<Monkey*> monkeys_;
};

class Monkey {
 public:
  Monkey(int id, MonkeyBarrel* others, std::function<int(int)> op,
         std::function<bool(int)> test, int succeed, int fail,
         std::vector<int> starting_items)
      : id_(id),
        others_(others),
        op_(op),
        test_(test),
        succeed_(succeed),
        fail_(fail),
        items_(starting_items) {}

  void TakeItem(int item) { items_.push_back(item); }
  int Activity() const { return total_actions_; }
  int Id() const { return id_; }

  void Play() {
    for (int i : items_) {
      l("playing with ", i);
      i = op_(i);
      i = i / 3;
      if (test_(i)) {
        others_->GetMonkey(succeed_)->TakeItem(i);
        l("throwing ", i, " to ", succeed_);
      } else {
        others_->GetMonkey(fail_)->TakeItem(i);
        l("throwing ", i, " to ", fail_);
      }
    }
    total_actions_ += items_.size();
    items_.clear();
    l("done turn!");
  }

 private:
  int id_;
  MonkeyBarrel* others_;
  std::function<int(int)> op_;
  std::function<bool(int)> test_;
  int succeed_;
  int fail_;

  int total_actions_ = 0;

  std::vector<int> items_;
};

int main(int argc, char** argv) {
  MonkeyBarrel monkeys;
  std::vector<std::unique_ptr<Monkey>> monkey_list;
  monkey_list.emplace_back(absl::make_unique<Monkey>(
      0,  //
      &monkeys, [](int i) -> int { return i * 17; },
      [](int i) -> bool { return i % 11 == 0; },  //
      2, 3,                                       //
      std::vector<int>{56, 52, 58, 96, 70, 75, 72}));
  monkey_list.emplace_back(absl::make_unique<Monkey>(
      1,  //
      &monkeys, [](int i) -> int { return i + 7; },
      [](int i) -> bool { return i % 3 == 0; },  //
      6, 5,                                      //
      std::vector<int>{75, 58, 86, 80, 55, 81}));
  monkey_list.emplace_back(absl::make_unique<Monkey>(
      2,  //
      &monkeys, [](int i) -> int { return i * i; },
      [](int i) -> bool { return i % 5 == 0; },  //
      1, 7,                                      //
      std::vector<int>{73, 68, 73, 90}));
  monkey_list.emplace_back(absl::make_unique<Monkey>(
      3,  //
      &monkeys, [](int i) -> int { return i + 1; },
      [](int i) -> bool { return i % 7 == 0; },  //
      2, 7,                                      //
      std::vector<int>{72, 89, 55, 51, 59}));
  monkey_list.emplace_back(absl::make_unique<Monkey>(
      4,  //
      &monkeys, [](int i) -> int { return i * 3; },
      [](int i) -> bool { return i % 19 == 0; },  //
      0, 3,                                       //
      std::vector<int>{76, 76, 91}));
  monkey_list.emplace_back(absl::make_unique<Monkey>(
      5,  //
      &monkeys, [](int i) -> int { return i + 4; },
      [](int i) -> bool { return i % 2 == 0; },  //
      6, 4,                                      //
      std::vector<int>{88}));
  monkey_list.emplace_back(absl::make_unique<Monkey>(
      6,  //
      &monkeys, [](int i) -> int { return i + 8; },
      [](int i) -> bool { return i % 13 == 0; },  //
      4, 0,                                       //
      std::vector<int>{64, 63, 56, 50, 77, 55, 55, 86}));
  monkey_list.emplace_back(absl::make_unique<Monkey>(
      7,  //
      &monkeys, [](int i) -> int { return i + 6; },
      [](int i) -> bool { return i % 17 == 0; },  //
      1, 5,                                       //
      std::vector<int>{79, 58}));
  for (auto& m : monkey_list) {
    monkeys.AddMonkey(m.get());
  }

  for (int round = 0; round < 20; round++) {
    l("====ROUND ", round, "======");
    for (auto& m : monkey_list) {
      l("MONKEY ", m->Id());
      m->Play();
    }
  }

  for (auto& m : monkey_list) {
    l("monkey ", m->Id(), " made ", m->Activity(), " actions");
  }
  std::sort(monkey_list.begin(), monkey_list.end(),
            [](const auto& lhs, const auto& rhs) {
              return lhs->Activity() < rhs->Activity();
            });

  l("most active monkey: ", monkey_list[7]->Id(), " with ",
    monkey_list[7]->Activity(), " actions.");
  l("second most active monkey: ", monkey_list[6]->Id(), " with ",
    monkey_list[6]->Activity(), " actions.");
  std::cout << "part 1: "
            << (monkey_list[7]->Activity() * monkey_list[6]->Activity())
            << std::endl;
}
