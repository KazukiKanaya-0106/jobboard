import { z } from "zod";
import { apiRequest } from "../../lib/api-client";
import type { StoredAuth } from "../../lib/storage";
import type { NodeDto, CreateNodeRequest } from "./schemas";
import { nodeSchema } from "./schemas";

export type Node = {
  id: number;
  nodeName: string;
  currentJobId: number | null;
  createdAt: Date;
};

const nodesArraySchema = z.array(nodeSchema);

function mapNode(dto: NodeDto): Node {
  return {
    id: dto.id,
    nodeName: dto.node_name,
    currentJobId: dto.current_job_id ?? null,
    createdAt: dto.created_at,
  };
}

export async function fetchNodes(auth: StoredAuth): Promise<Node[]> {
  const dto = await apiRequest(`/api/nodes`, {
    method: "GET",
    token: auth.token,
  });
  const nodes = nodesArraySchema.parse(dto);
  return nodes.map(mapNode);
}

type CreateNodeResponse = {
  node: Node;
  token?: string;
};

export async function createNode(auth: StoredAuth, request: CreateNodeRequest): Promise<CreateNodeResponse> {
  const dto = await apiRequest(`/api/nodes`, {
    method: "POST",
    token: auth.token,
    body: {
      node_name: request.nodeName,
    },
  });
  const parsed = nodeSchema.parse(dto);
  return {
    node: mapNode(parsed),
    token: parsed.node_token,
  };
}

export async function deleteNode(auth: StoredAuth, nodeId: number): Promise<void> {
  await apiRequest(`/api/nodes/${nodeId}`, {
    method: "DELETE",
    token: auth.token,
  });
}
