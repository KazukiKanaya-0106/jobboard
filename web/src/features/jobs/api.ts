import { apiRequest } from "../../lib/apiCient";
import type { StoredAuth } from "../../lib/storage";
import type { JobDto } from "./schemas";
import { jobsArraySchema } from "./schemas";

export type Job = {
  id: number;
  clusterId: string;
  nodeId: number;
  status: string;
  startedAt: Date;
  finishedAt: Date | null | undefined;
  durationHours: number | null | undefined;
  tag: string | null | undefined;
};

function parseDuration(value: JobDto["duration_hours"]): number | null {
  if (value === null || value === undefined) return null;
  if (typeof value === "number" && Number.isFinite(value)) {
    return value;
  }
  if (typeof value === "string") {
    const numeric = Number(value);
    return Number.isFinite(numeric) ? numeric : null;
  }
  if (typeof value === "object" && value && "Microseconds" in value) {
    const micro = (value as { Microseconds: number }).Microseconds;
    if (typeof micro === "number" && Number.isFinite(micro)) {
      return micro / (1000 * 1000) / 3600;
    }
  }
  return null;
}

function mapJob(dto: JobDto): Job {
  return {
    id: dto.id,
    clusterId: dto.cluster_id,
    nodeId: dto.node_id,
    status: dto.status,
    startedAt: dto.started_at,
    finishedAt: dto.finished_at,
    durationHours: parseDuration(dto.duration_hours),
    tag: dto.tag ?? null,
  };
}

export async function fetchJobs(auth: StoredAuth): Promise<Job[]> {
  const dto = await apiRequest(`/api/jobs`, {
    method: "GET",
    token: auth.token,
  });
  console.log("Raw job DTOs:", dto);
  const jobs = jobsArraySchema.parse(dto);
  console.log("Fetched jobs:", jobs);
  return jobs.map(mapJob);
}
