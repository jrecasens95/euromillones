type NumberBallProps = {
  value: number;
  muted?: boolean;
};

export function NumberBall({ value, muted }: NumberBallProps) {
  return (
    <span
      className={[
        "inline-flex h-9 w-9 shrink-0 items-center justify-center rounded-full text-sm font-extrabold shadow-sm",
        muted ? "bg-slate-200 text-slate-700" : "bg-white text-slate-950 ring-2 ring-cyan-300"
      ].join(" ")}
    >
      {value}
    </span>
  );
}
