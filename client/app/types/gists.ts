export interface Gist {
  fileId: string
  userId: string
  fileName: string

  gistTitle: string
  forkedFrom: string
  shortUrl: string

  viewCount: number
  isPublic: boolean
  isDeleted: boolean

  auditLog: string[]

  createdAt: string
  updatedAt: string

}
