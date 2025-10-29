import { z } from "zod";

export const nodeSchema = z.object({
  id: z.number(),
  node_name: z.string(),
  current_job_id: z.number().nullable().optional(),
  created_at: z.coerce.date(),
  node_token: z.string().optional(),
});

export type NodeDto = z.infer<typeof nodeSchema>;

export const createNodeRequestSchema = z.object({
  nodeName: z.string().min(1, "ノード名を入力してください").max(255, "ノード名は255文字以内で入力してください"),
});

export type CreateNodeRequest = z.infer<typeof createNodeRequestSchema>;
