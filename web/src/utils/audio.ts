export const formatSampleRate = (rate?: number | null) => {
  if (!rate) return null;
  return `${rate / 1000} KHz`;
};
