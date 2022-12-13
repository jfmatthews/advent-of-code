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

template <typename... Ts>
int DIE(Ts... ts) {
  (std::cerr << ... << ts) << std::endl;
  return 1;
}

class Packet {
 public:
  static Packet Parse(absl::string_view str) {
    if (str[0] == '[') {
      return Packet(parseList(str.substr(1, str.length() - 2)));
    } else {
      int val;
      auto ok = absl::SimpleAtoi(str, &val);
      if (ok) {
        l("parsed int packet ", val);
      } else {
        l("failed to parse ", str);
      }
      return Packet(val);
    }
  }

  Packet(const Packet&) = default;
  Packet(Packet&&) = default;
  Packet& operator=(const Packet&) = default;
  Packet& operator=(Packet&&) = default;

  friend std::ostream& operator<<(std::ostream&, const Packet&);

  auto operator<=>(const Packet& other) const {
    l("comparing ", *this, " and ", other);
    switch (type_) {
      case Type::SINGLE:
        switch (other.type_) {
          case Type::SINGLE:
            return val_ <=> other.val_;
          case Type::LIST:
            // convert this singular packet to a list, and continue
            return Packet(std::vector<Packet>{Packet(val_)}) <=> other;
        }
        break;
      case Type::LIST:
        switch (other.type_) {
          case Type::SINGLE:
            // convert that singular packet to a list, and continue
            return (*this) <=> Packet(std::vector<Packet>{Packet(other.val_)});
          case Type::LIST:
            for (int i = 0;
                 i < std::min(children_.size(), other.children_.size()); i++) {
              if (auto cmp = this->children_[i] <=> other.children_[i];
                  cmp != 0) {
                return cmp;
              }
            }
            return children_.size() <=> other.children_.size();
        }
    }
  }
  bool operator==(const Packet& other) const {
    if (type_ != other.type_) return false;
    return (*this <=> other) == 0;
  }

 private:
  explicit Packet(int val) : type_(Type::SINGLE), val_(val) {}
  explicit Packet(std::vector<Packet> nested)
      : type_(Type::LIST), children_(nested) {}

  static std::vector<Packet> parseList(absl::string_view str) {
    std::vector<Packet> parts;

    while (str.length() > 0) {
      l("continuing packet parse: ", str);
      if (str[0] != '[') {
        auto end_of_int = str.find_first_of(',');
        // l("parsing int packet: ", str.substr(0, end_of_int));
        parts.push_back(Packet::Parse(str.substr(0, end_of_int)));
        if (end_of_int == std::string::npos) {
          str = "";
        } else {
          str = str.substr(end_of_int + 1);
        }
        continue;
      }

      int depth = 1;
      auto end_of_list = std::string::npos;
      for (int i = 1; i < str.length(); i++) {
        if (str[i] == '[') depth++;
        if (str[i] == ']') depth--;
        if (depth == 0) {
          end_of_list = i;
          break;
        }
      }
      l("list string is ", str.substr(1, end_of_list - 1), "; descending");
      parts.push_back(Packet(parseList(str.substr(1, end_of_list - 1))));

      auto start_of_next_el = str.find(',', end_of_list);
      if (start_of_next_el == std::string::npos) {
        str = "";
      } else {
        str = str.substr(start_of_next_el + 1);
      }
    }

    l("parsed list of length ", parts.size());
    return parts;
  }

  enum class Type { SINGLE, LIST };
  Type type_;
  int val_;
  std::vector<Packet> children_;
};

std::ostream& operator<<(std::ostream& os, const Packet& p) {
  switch (p.type_) {
    case Packet::Type::SINGLE:
      return os << p.val_;
    case Packet::Type::LIST:
      os << "[";
      for (int i = 0; i < p.children_.size(); i++) {
        if (i > 0) {
          os << ",";
        }
        os << p.children_[i];
      }
      return os << "]";
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

  std::vector<Packet> packets;
  for (std::string line; std::getline(input, line);) {
    if (line == "") {
      continue;
    }
    packets.push_back(Packet::Parse(line));
  }

  int total = 0;
  for (int i = 1; i < packets.size(); i += 2) {
    const int pair_num = (i / 2) + 1;
    if (packets[i - 1] <= packets[i]) {
      l("PAIR ", pair_num, " MATCHES");
      total += pair_num;
    } else {
      l("PAIR ", pair_num, " DOES NOT MATCH");
    }
  }
  std::cout << "part 1: " << total << std::endl;

  Packet packet2 = Packet::Parse("[[2]]");
  Packet packet6 = Packet::Parse("[[6]]");
  packets.push_back(packet2);
  packets.push_back(packet6);
  std::sort(packets.begin(), packets.end());

  std::cout << "====SORTED====" << std::endl;

  int score = 1;
  for (int i = 0; i < packets.size(); i++) {
    std::cout << packets[i] << std::endl;
    if (packets[i] == packet2 || packets[i] == packet6) {
      score *= (i + 1);
    }
  }
  std::cout << "part 2: " << score << std::endl;

  return 0;
}
