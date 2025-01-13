export type ApiResponse<T> = {
  data: T | null;
  err_code: number | null;
};
