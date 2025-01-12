// eslint-disable-next-line ts/consistent-type-definitions
export type Gist = {
  fieldId: string
  userId: string

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
