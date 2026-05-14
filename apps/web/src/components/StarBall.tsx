type StarBallProps = {
  value: number;
};

export function StarBall({ value }: StarBallProps) {
  return (
    <span className="inline-flex h-9 w-9 shrink-0 items-center justify-center rounded-full bg-amber-300 text-sm font-extrabold text-slate-950 shadow-sm ring-2 ring-amber-100">
      {value}
    </span>
  );
}
