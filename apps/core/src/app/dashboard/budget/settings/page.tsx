"use client";

import { motion } from "motion/react";
import { useState } from "react";
import {
  useCategories,
  useTags,
  useCreateCategory,
  useUpdateCategory,
  useDeleteCategory,
  useCreateTag,
  useUpdateTag,
  useDeleteTag,
} from "@/hooks/use-budget";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Plus, Trash2, Edit2, Save, X, ArrowLeft } from "lucide-react";
import { Badge } from "@/components/ui/badge";
import { CategoryTypeSelector, CategoryTypeBadge } from "@/components/budget";
import type { Category, Tag } from "@/types/api";
import Link from "next/link";
import { useToast } from "@/components/ui/toast";
import { GuestBlockedError } from "@/hooks/use-budget";

export default function SettingsPage() {
  const { data: categories, isLoading: categoriesLoading } = useCategories();
  const { data: tags, isLoading: tagsLoading } = useTags();

  const createCategory = useCreateCategory();
  const updateCategory = useUpdateCategory();
  const deleteCategory = useDeleteCategory();
  const createTag = useCreateTag();
  const updateTag = useUpdateTag();
  const deleteTag = useDeleteTag();
  const { showToast } = useToast();

  // Category form state
  const [newCategory, setNewCategory] = useState({ name: "", color: "#3b82f6", icon: "" });
  const [editingCategory, setEditingCategory] = useState<Category | null>(null);

  // Tag form state
  const [newTag, setNewTag] = useState({ name: "", color: "#8b5cf6" });
  const [editingTag, setEditingTag] = useState<Tag | null>(null);

  // Category handlers
  const handleCreateCategory = async () => {
    if (!newCategory.name.trim()) return;
    try {
      await createCategory.mutateAsync(newCategory);
      showToast("Category created successfully", "success");
      setNewCategory({ name: "", color: "#3b82f6", icon: "" });
    } catch (error) {
      if (error instanceof GuestBlockedError) {
        showToast(error.message, "warning");
      } else {
        showToast("Failed to create category", "error");
      }
    }
  };

  const handleUpdateCategory = async () => {
    if (!editingCategory) return;
    try {
      await updateCategory.mutateAsync({
        id: editingCategory.id,
        data: {
          name: editingCategory.name,
          color: editingCategory.color,
          icon: editingCategory.icon,
        },
      });
      showToast("Category updated successfully", "success");
      setEditingCategory(null);
    } catch (error) {
      if (error instanceof GuestBlockedError) {
        showToast(error.message, "warning");
      } else {
        showToast("Failed to update category", "error");
      }
    }
  };

  const handleDeleteCategory = async (id: string) => {
    if (confirm("Delete this category? This will not delete expenses.")) {
      try {
        await deleteCategory.mutateAsync(id);
        showToast("Category deleted successfully", "success");
      } catch (error) {
        if (error instanceof GuestBlockedError) {
          showToast(error.message, "warning");
        } else {
          showToast("Failed to delete category", "error");
        }
      }
    }
  };

  // Tag handlers
  const handleCreateTag = async () => {
    if (!newTag.name.trim()) return;
    try {
      await createTag.mutateAsync(newTag);
      showToast("Tag created successfully", "success");
      setNewTag({ name: "", color: "#8b5cf6" });
    } catch (error) {
      if (error instanceof GuestBlockedError) {
        showToast(error.message, "warning");
      } else {
        showToast("Failed to create tag", "error");
      }
    }
  };

  const handleUpdateTag = async () => {
    if (!editingTag) return;
    try {
      await updateTag.mutateAsync({
        id: editingTag.id,
        data: {
          name: editingTag.name,
          color: editingTag.color,
        },
      });
      showToast("Tag updated successfully", "success");
      setEditingTag(null);
    } catch (error) {
      if (error instanceof GuestBlockedError) {
        showToast(error.message, "warning");
      } else {
        showToast("Failed to update tag", "error");
      }
    }
  };

  const handleDeleteTag = async (id: string) => {
    if (confirm("Delete this tag? This will remove it from all expenses.")) {
      try {
        await deleteTag.mutateAsync(id);
        showToast("Tag deleted successfully", "success");
      } catch (error) {
        if (error instanceof GuestBlockedError) {
          showToast(error.message, "warning");
        } else {
          showToast("Failed to delete tag", "error");
        }
      }
    }
  };

  return (
      <div className="px-4 md:px-6 py-6">
      {/* Header */}
      <motion.div
        className="mb-8"
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.5 }}
      >
        <Link href="/dashboard/budget">
          <Button variant="ghost" size="sm" className="mb-4">
            <ArrowLeft className="mr-2 h-4 w-4" />
            Back to Budget
          </Button>
        </Link>
        <h1 className="text-primary mb-2">BUDGET SETTINGS</h1>
        <p className="text-sm text-muted-foreground">
          Manage your categories and tags
        </p>
      </motion.div>

      <div className="grid gap-8 lg:grid-cols-2">
        {/* Categories Section */}
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.5, delay: 0.1 }}
        >
          <Card>
            <CardHeader>
              <CardTitle>Categories</CardTitle>
              <CardDescription>Organize expenses by category</CardDescription>
            </CardHeader>
            <CardContent className="space-y-4">
              {/* Add Category Form */}
              <div className="space-y-3 p-4 border rounded-lg">
                <Label>Add New Category</Label>
                <div className="grid gap-2">
                  <Input
                    placeholder="Category name"
                    value={newCategory.name}
                    onChange={(e) => setNewCategory({ ...newCategory, name: e.target.value })}
                  />
                  <div className="grid grid-cols-2 gap-2">
                    <Input
                      placeholder="Icon (emoji)"
                      value={newCategory.icon}
                      onChange={(e) => setNewCategory({ ...newCategory, icon: e.target.value })}
                    />
                    <Input
                      type="color"
                      value={newCategory.color}
                      onChange={(e) => setNewCategory({ ...newCategory, color: e.target.value })}
                    />
                  </div>
                  <Button
                    onClick={handleCreateCategory}
                    disabled={!newCategory.name.trim() || createCategory.isPending}
                    size="sm"
                  >
                    <Plus className="mr-2 h-4 w-4" />
                    Add Category
                  </Button>
                </div>
              </div>

              {/* Categories List */}
              <div className="space-y-2">
                {categoriesLoading ? (
                  <p className="text-sm text-muted-foreground">Loading...</p>
                ) : categories && categories.length > 0 ? (
                  categories.map((category) => (
                    <div
                      key={category.id}
                      className="flex items-center justify-between p-3 border rounded-lg"
                    >
                      {editingCategory?.id === category.id ? (
                        <div className="flex-1 grid gap-2">
                          <Input
                            value={editingCategory.name}
                            onChange={(e) =>
                              setEditingCategory({ ...editingCategory, name: e.target.value })
                            }
                          />
                          <div className="grid grid-cols-2 gap-2">
                            <Input
                              value={editingCategory.icon || ""}
                              onChange={(e) =>
                                setEditingCategory({ ...editingCategory, icon: e.target.value })
                              }
                            />
                            <Input
                              type="color"
                              value={editingCategory.color || "#3b82f6"}
                              onChange={(e) =>
                                setEditingCategory({ ...editingCategory, color: e.target.value })
                              }
                            />
                          </div>
                          <div className="flex gap-2">
                            <Button size="sm" onClick={handleUpdateCategory}>
                              <Save className="mr-2 h-4 w-4" />
                              Save
                            </Button>
                            <Button
                              size="sm"
                              variant="outline"
                              onClick={() => setEditingCategory(null)}
                            >
                              <X className="mr-2 h-4 w-4" />
                              Cancel
                            </Button>
                          </div>
                        </div>
                      ) : (
                        <>
                          <div className="flex items-center gap-2">
                            <Badge
                              style={{
                                backgroundColor: category.color ? `${category.color}20` : undefined,
                                borderColor: category.color || undefined,
                                color: category.color || "inherit",
                              }}
                            >
                              {category.icon && <span className="mr-1">{category.icon}</span>}
                              {category.name}
                            </Badge>
                            <CategoryTypeBadge type={category.category_type} />
                          </div>
                          <div className="flex items-center gap-2">
                            <CategoryTypeSelector category={category} />
                            <Button
                              size="sm"
                              variant="ghost"
                              onClick={() => setEditingCategory(category)}
                            >
                              <Edit2 className="h-4 w-4" />
                            </Button>
                            <Button
                              size="sm"
                              variant="ghost"
                              className="text-destructive"
                              onClick={() => handleDeleteCategory(category.id)}
                            >
                              <Trash2 className="h-4 w-4" />
                            </Button>
                          </div>
                        </>
                      )}
                    </div>
                  ))
                ) : (
                  <p className="text-sm text-muted-foreground text-center py-4">
                    No categories yet
                  </p>
                )}
              </div>
            </CardContent>
          </Card>
        </motion.div>

        {/* Tags Section */}
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.5, delay: 0.2 }}
        >
          <Card>
            <CardHeader>
              <CardTitle>Tags</CardTitle>
              <CardDescription>Add labels to your expenses</CardDescription>
            </CardHeader>
            <CardContent className="space-y-4">
              {/* Add Tag Form */}
              <div className="space-y-3 p-4 border rounded-lg">
                <Label>Add New Tag</Label>
                <div className="grid gap-2">
                  <Input
                    placeholder="Tag name"
                    value={newTag.name}
                    onChange={(e) => setNewTag({ ...newTag, name: e.target.value })}
                  />
                  <Input
                    type="color"
                    value={newTag.color}
                    onChange={(e) => setNewTag({ ...newTag, color: e.target.value })}
                  />
                  <Button
                    onClick={handleCreateTag}
                    disabled={!newTag.name.trim() || createTag.isPending}
                    size="sm"
                  >
                    <Plus className="mr-2 h-4 w-4" />
                    Add Tag
                  </Button>
                </div>
              </div>

              {/* Tags List */}
              <div className="space-y-2">
                {tagsLoading ? (
                  <p className="text-sm text-muted-foreground">Loading...</p>
                ) : tags && tags.length > 0 ? (
                  tags.map((tag) => (
                    <div
                      key={tag.id}
                      className="flex items-center justify-between p-3 border rounded-lg"
                    >
                      {editingTag?.id === tag.id ? (
                        <div className="flex-1 grid gap-2">
                          <Input
                            value={editingTag.name}
                            onChange={(e) => setEditingTag({ ...editingTag, name: e.target.value })}
                          />
                          <Input
                            type="color"
                            value={editingTag.color || "#8b5cf6"}
                            onChange={(e) => setEditingTag({ ...editingTag, color: e.target.value })}
                          />
                          <div className="flex gap-2">
                            <Button size="sm" onClick={handleUpdateTag}>
                              <Save className="mr-2 h-4 w-4" />
                              Save
                            </Button>
                            <Button
                              size="sm"
                              variant="outline"
                              onClick={() => setEditingTag(null)}
                            >
                              <X className="mr-2 h-4 w-4" />
                              Cancel
                            </Button>
                          </div>
                        </div>
                      ) : (
                        <>
                          <Badge
                            style={{
                              backgroundColor: tag.color ? `${tag.color}20` : undefined,
                              borderColor: tag.color || undefined,
                              color: tag.color || "inherit",
                            }}
                          >
                            {tag.name}
                          </Badge>
                          <div className="flex gap-1">
                            <Button
                              size="sm"
                              variant="ghost"
                              onClick={() => setEditingTag(tag)}
                            >
                              <Edit2 className="h-4 w-4" />
                            </Button>
                            <Button
                              size="sm"
                              variant="ghost"
                              className="text-destructive"
                              onClick={() => handleDeleteTag(tag.id)}
                            >
                              <Trash2 className="h-4 w-4" />
                            </Button>
                          </div>
                        </>
                      )}
                    </div>
                  ))
                ) : (
                  <p className="text-sm text-muted-foreground text-center py-4">
                    No tags yet
                  </p>
                )}
              </div>
            </CardContent>
          </Card>
        </motion.div>
      </div>
    </div>
  );
}
