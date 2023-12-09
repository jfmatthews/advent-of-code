#ifndef LIB_LOGGING_H_
#define LIB_LOGGING_H_

namespace aoc {

template <typename... Ts>
void l(Ts... ts) {
#ifndef NDEBUG
  (std::cout << ... << ts) << std::endl;
#endif
}

template <typename... Ts>
int DIE(Ts... ts) {
  (std::cerr << ... << ts) << std::endl;
  return 1;
}

}  // namespace aoc

#endif
