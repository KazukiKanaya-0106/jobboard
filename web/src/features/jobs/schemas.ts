import { z } from "zod";

const intervalObjectSchema = z
  .object({
    Microseconds: z.number(),
  })
  .passthrough();

const durationSchema = z.union([z.number(), z.string(), intervalObjectSchema, z.null(), z.undefined()]);

export const jobSchema = z.object({
  id: z.number(),
  cluster_id: z.string(),
  node_id: z.number(),
  started_at: z.coerce.date(),
  finished_at: z.coerce.date().nullable().optional(),
  status: z.string(),
  duration_hours: durationSchema.optional(),
  tag: z.string().nullable().optional(),
  error_text: z.string().nullable().optional(),
});

export type JobDto = z.infer<typeof jobSchema>;

export const jobsArraySchema = z.array(jobSchema);
export type JobsArrayDto = z.infer<typeof jobsArraySchema>;
