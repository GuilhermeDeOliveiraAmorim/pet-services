import { isAxiosError } from "axios";
import type { ProblemDetails, ProblemDetailsResponse } from "@/application";

export const getApiProblemDetails = (
  error: unknown,
): ProblemDetails | undefined => {
  if (!isAxiosError<ProblemDetailsResponse>(error)) {
    return undefined;
  }

  return error.response?.data?.errors?.[0];
};

export const getApiErrorMessage = (
  error: unknown,
  fallbackMessage: string,
): string => {
  const problem = getApiProblemDetails(error);

  if (problem?.detail) {
    return problem.detail;
  }

  if (problem?.title) {
    return problem.title;
  }

  if (error instanceof Error && error.message.trim()) {
    return error.message;
  }

  return fallbackMessage;
};

export const isUnverifiedEmailError = (error: unknown): boolean => {
  const problem = getApiProblemDetails(error);

  if (!problem) {
    return false;
  }

  const title = problem.title?.toLowerCase() ?? "";
  const detail = problem.detail?.toLowerCase() ?? "";

  return (
    problem.status === 403 &&
    (title.includes("email não verificado") ||
      detail.includes("email") ||
      detail.includes("verifica"))
  );
};
