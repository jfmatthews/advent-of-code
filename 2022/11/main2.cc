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
constexpr const int64_t kLcm = 9699690;

template <typename... Ts>
void l(Ts... ts) {
  if (debug) {
    (std::cout << ... << ts) << std::endl;
  }
}

template <typename... Ts>
int64_t DIE(Ts... ts) {
  (std::cerr << ... << ts) << std::endl;
  return 1;
}

class Monkey;

class MonkeyBarrel {
 public:
  Monkey* GetMonkey(int64_t idx) const {
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
  Monkey(int64_t id, MonkeyBarrel* others, std::function<int64_t(int64_t)> op,
         int64_t mod, int64_t succeed, int64_t fail,
         std::vector<int64_t> starting_items)
      : id_(id),
        others_(others),
        op_(op),
        mod_(mod),
        succeed_(succeed),
        fail_(fail),
        items_(starting_items) {}

  void TakeItem(int64_t item) { items_.push_back(item); }
  int64_t Activity() const { return total_actions_; }
  int64_t Id() const { return id_; }

  void Play() {
    for (int64_t i : items_) {
      l("playing with ", i);
      i = op_(i);
      i %= kLcm;
      if (i % mod_ == 0) {
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
  int64_t id_;
  MonkeyBarrel* others_;
  std::function<int64_t(int64_t)> op_;
  int64_t mod_;
  int64_t succeed_;
  int64_t fail_;

  int64_t total_actions_ = 0;

  std::vector<int64_t> items_;
};

int main(int argc, char** argv) {
  MonkeyBarrel monkeys;
  std::vector<std::unique_ptr<Monkey>> monkey_list;
  monkey_list.emplace_back(absl::make_unique<Monkey>(
      0,  //
      &monkeys, [](int64_t i) -> int64_t { return i * 17; },
      11,    //
      2, 3,  //
      std::vector<int64_t>{56, 52, 58, 96, 70, 75, 72}));
  monkey_list.emplace_back(absl::make_unique<Monkey>(
      1,  //
      &monkeys, [](int64_t i) -> int64_t { return i + 7; },
      3,     //
      6, 5,  //
      std::vector<int64_t>{75, 58, 86, 80, 55, 81}));
  monkey_list.emplace_back(absl::make_unique<Monkey>(
      2,  //
      &monkeys, [](int64_t i) -> int64_t { return i * i; },
      5,     //
      1, 7,  //
      std::vector<int64_t>{73, 68, 73, 90}));
  monkey_list.emplace_back(absl::make_unique<Monkey>(
      3,  //
      &monkeys, [](int64_t i) -> int64_t { return i + 1; },
      7,     //
      2, 7,  //
      std::vector<int64_t>{72, 89, 55, 51, 59}));
  monkey_list.emplace_back(absl::make_unique<Monkey>(
      4,  //
      &monkeys, [](int64_t i) -> int64_t { return i * 3; },
      19,    //
      0, 3,  //
      std::vector<int64_t>{76, 76, 91}));
  monkey_list.emplace_back(absl::make_unique<Monkey>(
      5,  //
      &monkeys, [](int64_t i) -> int64_t { return i + 4; },
      2,     //
      6, 4,  //
      std::vector<int64_t>{88}));
  monkey_list.emplace_back(absl::make_unique<Monkey>(
      6,  //
      &monkeys, [](int64_t i) -> int64_t { return i + 8; },
      13,    //
      4, 0,  //
      std::vector<int64_t>{64, 63, 56, 50, 77, 55, 55, 86}));
  monkey_list.emplace_back(absl::make_unique<Monkey>(
      7,  //
      &monkeys, [](int64_t i) -> int64_t { return i + 6; },
      17,    //
      1, 5,  //
      std::vector<int64_t>{79, 58}));
  for (auto& m : monkey_list) {
    monkeys.AddMonkey(m.get());
  }

  for (int64_t round = 0; round < 10000; round++) {
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
  std::cout << "part 2: "
            << (monkey_list[7]->Activity() * monkey_list[6]->Activity())
            << std::endl;
}
