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

constexpr const int kBigSizeThreshold = 100000;
constexpr const int kDiskSize = 70000000;
constexpr const int kDiskHeadroom = 30000000;

class INode {
 public:
  static std::unique_ptr<INode> Dir(INode* parent, std::string name) {
    return absl::WrapUnique(new INode(parent, name));
  }
  static std::unique_ptr<INode> File(INode* parent, std::string name,
                                     int size) {
    return absl::WrapUnique(new INode(parent, name, size));
  }

  absl::string_view name() const { return name_; }
  std::string full_path() const {
    if (parent_ == nullptr) {
      return name_;
    } else {
      return absl::StrCat(parent_->full_path(), "/", name_);
    }
  }
  bool is_dir() const { return type_ == Type::DIR; }
  int GetTotalSize() const;
  INode* GetChild(absl::string_view name) const {
    for (auto* c : children_) {
      if (c->name_ == name) {
        return c;
      }
    }
    return nullptr;
  }
  INode* parent() const { return parent_; }

  bool AddChild(INode* other) {
    if (type_ != Type::DIR) {
      throw absl::StrCat("adding child to non-directory ", name_);
    }
    if (other->parent_ != this) {
      throw absl::StrCat("trying to add ", other->name_, " to ", name_,
                         " but its parent is ", other->parent_->name_);
    }
    if (GetChild(other->name_) == nullptr) {
      children_.push_back(other);
      return true;
    } else {
      return false;
    }
  }

 private:
  explicit INode(INode* parent, std::string name)
      : type_(Type::DIR), name_(name), parent_(parent) {}
  explicit INode(INode* parent, std::string name, int fsize)
      : type_(Type::FILE), name_(name), parent_(parent), size_(fsize) {}

  enum class Type { DIR = 0, FILE };

  Type type_;
  std::string name_;
  INode* parent_ = nullptr;
  std::vector<INode*> children_;
  int size_;
};

int INode::GetTotalSize() const {
  switch (type_) {
    case Type::FILE:
      return size_;
    case Type::DIR:
      int total = 0;
      for (auto& c : children_) {
        total += c->GetTotalSize();
      }
      return total;
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

  {
    std::ifstream input(abs_inpath.string());

    std::vector<std::unique_ptr<INode>> nodes;
    nodes.emplace_back(INode::Dir(nullptr, "/"));
    INode* root = nodes[0].get();
    INode* pwd = root;

    bool in_ls = false;
    for (std::string line; std::getline(input, line);) {
      if (line == "$ cd /") {
        in_ls = false;
        pwd = root;
        continue;
      }
      std::vector<absl::string_view> tokens = absl::StrSplit(line, " ");
      if (tokens[0] == "$") {
        in_ls = false;
      } else {
        if (!in_ls) {
          std::cerr << "line didn't start with $ prompt: " << line << std::endl;
          return 1;
        }

        if (tokens.size() != 2) {
          std::cerr << "bad entry line " << line << std::endl;
          return 1;
        }

        std::unique_ptr<INode> new_child;
        if (tokens[0] == "dir") {
          new_child = INode::Dir(pwd, std::string(tokens[1]));
        } else {
          int size;
          bool ok = absl::SimpleAtoi(tokens[0], &size);
          if (!ok) {
            std::cerr << "bad entry line (can't parse size): " << line
                      << std::endl;
            return 1;
          }
          new_child = INode::File(pwd, std::string(tokens[1]), size);
        }
        if (pwd->AddChild(new_child.get())) {
          nodes.emplace_back(std::move(new_child));
        }
      }

      if (tokens[1] == "cd") {
        if (tokens.size() != 3) {
          std::cerr << "bad entry line " << line << std::endl;
          return 1;
        }

        if (tokens[2] == "..") {
          pwd = pwd->parent();
        } else {
          pwd = pwd->GetChild(tokens[2]);
        }
      }

      if (tokens[1] == "ls") {
        in_ls = true;
      }
    }

    // done parsing/running
    int total_big_sizes = 0;
    const int total_disk_used = nodes[0]->GetTotalSize();
    const int disk_to_free = total_disk_used - (kDiskSize - kDiskHeadroom);
    int best_delete_option = std::numeric_limits<int>::max();
    for (auto& n : nodes) {
      if (!n->is_dir()) continue;

      const int size = n->GetTotalSize();
      if (size < kBigSizeThreshold) {
        total_big_sizes += size;
        std::cout << n->full_path() << " has total size " << size
                  << " and is SMALL" << std::endl;
      } else {
        std::cout << n->full_path() << " has total size " << size
                  << " and is LARGE" << std::endl;
      }
      if (size >= disk_to_free && size < best_delete_option) {
        std::cout << "new best deletion candidate " << n->full_path()
                  << std::endl;
        best_delete_option = size;
      }
    }
    std::cout << "part 1: " << total_big_sizes << std::endl;
    std::cout << "part 2: " << best_delete_option << std::endl;
  }
}
